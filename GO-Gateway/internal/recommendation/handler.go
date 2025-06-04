package recommendation

import (
	"log"
	"strconv"
	"strings"

	"guru-game/internal/boardgame/service_board"
	"guru-game/internal/db/repository/user_states"

	"github.com/gofiber/fiber/v2"
)

// RecommendationClient defines the interface for recommendation service clients
type RecommendationClient interface {
	SendUserAction(action UserAction) error
	GetRecommendations(userID string, limit int) ([]Boardgame, error)
	SendAllBoardgames(boardgames []Boardgame) error
	GetAllBoardgames() ([]Boardgame, error)
	GetPopularBoardgames(limit int) ([]Boardgame, error)
	GetUserActions(userID string) ([]UserAction, error)
	GetBoardgameActions(boardgameID string) ([]UserAction, error)
}

// Handler handles recommendation-related HTTP requests
type Handler struct {
	client        RecommendationClient
	bgService     *service_board.BoardgameService
	userStateRepo user_states.UserStateRepository
}

// NewHandler creates a new recommendation handler
func NewHandler(client RecommendationClient, bgService *service_board.BoardgameService, userStateRepo user_states.UserStateRepository) *Handler {
	return &Handler{
		client:        client,
		bgService:     bgService,
		userStateRepo: userStateRepo,
	}
}

// HandleSendAllBoardgames handles sending all boardgames to the recommendation service
func (h *Handler) HandleSendAllBoardgames(c *fiber.Ctx) error {
	// Query ข้อมูลจาก DB
	boardgames, err := h.bgService.GetAllBoardgames()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get boardgames from database",
		})
	}

	// แปลงข้อมูลเป็น format ที่ Python service ต้องการ
	var recoBoardgames []Boardgame
	for _, bg := range boardgames {
		recoBoardgames = append(recoBoardgames, Boardgame{
			ID:              bg.ID,
			Title:           bg.Title,
			Description:     bg.Description,
			MinPlayers:      bg.MinPlayers,
			MaxPlayers:      bg.MaxPlayers,
			PlayTimeMin:     bg.PlayTimeMin,
			PlayTimeMax:     bg.PlayTimeMax,
			Categories:      bg.Categories,
			RatingAvg:       bg.RatingAvg,
			RatingCount:     bg.RatingCount,
			PopularityScore: bg.PopularityScore,
			ImageURL:        bg.ImageURL,
		})
	}

	// ส่งข้อมูลไปยัง Python service
	if err := h.client.SendAllBoardgames(recoBoardgames); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to send boardgames to recommendation service",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Successfully sent all boardgames to recommendation service",
		"count":   len(recoBoardgames),
	})
}

// HandleGetRecommendations handles getting recommendations for a user
func (h *Handler) HandleGetRecommendations(c *fiber.Ctx) error {
	userID := c.Query("user_id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id is required",
		})
	}

	limitStr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid limit parameter",
		})
	}

	recommendations, err := h.client.GetRecommendations(userID, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get recommendations",
		})
	}

	return c.JSON(fiber.Map{
		"boardgames": recommendations,
	})
}

// HandleGetPopularBoardgames handles getting popular boardgames
func (h *Handler) HandleGetPopularBoardgames(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid limit parameter",
		})
	}

	boardgames, err := h.client.GetPopularBoardgames(limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get popular boardgames",
		})
	}

	return c.JSON(fiber.Map{
		"boardgames": boardgames,
	})
}

// HandleAddUserAction handles adding a new user action
func (h *Handler) HandleAddUserAction(c *fiber.Ctx) error {
	var action UserAction
	if err := c.BodyParser(&action); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := h.client.SendUserAction(action); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to add user action",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User action added successfully",
	})
}

// HandleGetUserActions handles getting all actions for a user, and can filter by action type
func (h *Handler) HandleGetUserActions(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id is required",
		})
	}

	actions, err := h.client.GetUserActions(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get user actions",
		})
	}

	// Check if the route is for favorites and filter actions
	if strings.Contains(c.Path(), "/favorites/") {
		filteredActions := []UserAction{}
		for _, action := range actions {
			if action.ActionType == "favorite" {
				filteredActions = append(filteredActions, action)
			}
		}
		return c.JSON(fiber.Map{
			"actions": filteredActions,
		})
	}

	return c.JSON(fiber.Map{
		"actions": actions,
	})
}

// HandleGetBoardgameActions handles getting all actions for a boardgame
func (h *Handler) HandleGetBoardgameActions(c *fiber.Ctx) error {
	boardgameID := c.Params("boardgame_id")
	if boardgameID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "boardgame_id is required",
		})
	}

	actions, err := h.client.GetBoardgameActions(boardgameID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get boardgame actions",
		})
	}

	return c.JSON(fiber.Map{
		"actions": actions,
	})
}

// HandleGetAllBoardgamesFromES handles getting all boardgames from Elasticsearch via the recommendation service
func (h *Handler) HandleGetAllBoardgamesFromES(c *fiber.Ctx) error {
	boardgames, err := h.client.GetAllBoardgames()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get all boardgames from recommendation service",
		})
	}

	return c.JSON(fiber.Map{
		"boardgames": boardgames,
	})
}

// HandleGetFavoritedBoardgames handles fetching favorited boardgames for a user directly from DB
func (h *Handler) HandleGetFavoritedBoardgames(c *fiber.Ctx) error {
	log.Printf("Received request to get favorited boardgames for user ID: %s", c.Params("user_id"))

	userIDStr := c.Params("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user ID format",
		})
	}

	favoritedStates, err := h.userStateRepo.GetFavoritedByUserID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get favorited user states",
		})
	}

	var favoritedBoardgames []FavoritedBoardgame
	for _, state := range favoritedStates {
		boardgame, err := h.bgService.GetBoardGameByID(state.BoardgameID)
		if err != nil {
			// Log the error but continue processing other favorites
			log.Printf("Could not retrieve boardgame ID %d for user %d favorite: %v", state.BoardgameID, userID, err)
			continue // Skip this favorited item if boardgame details cannot be fetched
		}

		favoritedBoardgames = append(favoritedBoardgames, FavoritedBoardgame{
			UserID:          state.UserID,
			BoardgameID:     state.BoardgameID,
			Liked:           state.Liked,
			Favorited:       state.Favorited,
			Rating:          state.Rating,
			UpdatedAt:       state.UpdatedAt,
			Title:           boardgame.Title,
			Description:     boardgame.Description,
			MinPlayers:      boardgame.MinPlayers,
			MaxPlayers:      boardgame.MaxPlayers,
			PlayTimeMin:     boardgame.PlayTimeMin,
			PlayTimeMax:     boardgame.PlayTimeMax,
			Categories:      boardgame.Categories,
			RatingAvg:       boardgame.RatingAvg,
			RatingCount:     boardgame.RatingCount,
			PopularityScore: boardgame.PopularityScore,
			ImageURL:        boardgame.ImageURL,
		})
	}

	return c.JSON(fiber.Map{
		"favorites": favoritedBoardgames,
	})
}

package handlers

import (
	"log"
	"strconv"
	"time"

	"guru-game/internal/db/repository/user_states"

	"github.com/gofiber/fiber/v2"
)

// GameStateUpdateData represents the nested data structure for updates
type GameStateUpdateData map[string]interface{}

// GameStateUpdate represents the expected structure of the incoming request body
type GameStateUpdate struct {
	UserID     int                 `json:"user_id"` // Changed to int based on DB schema
	GameIDStr  string              `json:"game_id"` // Temporarily read as string
	UpdateData GameStateUpdateData `json:"state"`   // Changed json tag to "state"
}

// GameStateHandlers holds the necessary dependencies for game state handlers
type GameStateHandlers struct {
	UserStateRepo user_states.UserStateRepository
}

// NewGameStateHandlers creates a new GameStateHandlers instance
func NewGameStateHandlers(userStateRepo user_states.UserStateRepository) *GameStateHandlers {
	return &GameStateHandlers{
		UserStateRepo: userStateRepo,
	}
}

// HandleGameStateUpdate receives and processes game state updates
func (h *GameStateHandlers) HandleGameStateUpdate(c *fiber.Ctx) error {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("\n\n[%s] ===== GAME STATE UPDATE RECEIVED =====\n", timestamp)

	// Log raw request body
	rawBody := c.Body()
	log.Printf("Raw Request Body: %s", string(rawBody))

	var update GameStateUpdate

	// Parse the request body
	if err := c.BodyParser(&update); err != nil {
		log.Printf("[%s] Error parsing game state update request body: %v", timestamp, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to parse request body",
		})
	}

	// Convert game_id from string to int
	gameID, err := strconv.Atoi(update.GameIDStr)
	if err != nil {
		log.Printf("[%s] Error converting game_id to int: %v", timestamp, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid game_id format",
		})
	}

	// Extract data from updateData (state)
	// Safely extract boolean and float values with type assertions
	liked, _ := update.UpdateData["is_liked"].(bool)
	favorited, _ := update.UpdateData["is_favorite"].(bool)
	// Need to handle both int and float64 for rating if necessary, though float64 is safer
	rating, ok := update.UpdateData["userRating"].(float64)
	if !ok {
		// If it's not a float64, try int (e.g., if 0 was sent as 0 instead of 0.0)
		intRating, ok := update.UpdateData["userRating"].(int)
		if ok {
			rating = float64(intRating)
		} else {
			log.Printf("[%s] Warning: userRating is not a number type: %T", timestamp, update.UpdateData["userRating"])
			// Default rating to 0.0 or handle as an error if rating is mandatory
			rating = 0.0
		}
	}

	// Create UserState object
	userState := &user_states.UserState{
		UserID:      update.UserID,
		BoardgameID: gameID, // Use the converted int gameID
		Liked:       liked,
		Favorited:   favorited,
		Rating:      rating,
		UpdatedAt:   time.Now(),
	}

	// Save or update user state in the database
	err = h.UserStateRepo.SaveOrUpdate(c.Context(), userState)
	if err != nil {
		log.Printf("[%s] Error saving or updating user state: %v", timestamp, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to save or update user state",
		})
	}

	log.Printf("[%s] ===== GAME STATE UPDATE PROCESSED =====\n", timestamp)

	// Return the processed data in the desired format (using original string game_id for frontend consistency if needed, or converted int)
	// Returning the converted int game_id might be better for frontend consistency with type
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user_id": update.UserID,
		"game_id": gameID, // Return as int
		"state":   update.UpdateData,
	})
}

package handlers_board

import (
	"log"

	"guru-game/internal/boardgame/service_board"
	"guru-game/internal/db/repository/boardgame"

	"github.com/gofiber/fiber/v2"
)

// BoardGameHandlers struct holds dependencies for board game handlers
type BoardGameHandlers struct {
	BoardGameRepo boardgame.BoardGameRepository
}

// NewBoardGameHandlers creates a new BoardGameHandlers instance
func NewBoardGameHandlers(repo boardgame.BoardGameRepository) *BoardGameHandlers {
	return &BoardGameHandlers{
		BoardGameRepo: repo,
	}
}

// GetAllBoardGamesHandler handles fetching all board games
func (h *BoardGameHandlers) HandleGetAllBoardGames(c *fiber.Ctx) error {
	// Attempt to get userID from context locals
	userID, ok := c.Locals("userID").(int)
	if !ok {
		userID = 0 // Default to unauthenticated user if no userID is found
	}

	boardgames, err := service_board.GetAllBoardGames(userID, h.BoardGameRepo) // Pass userID and repo to service
	if err != nil {
		log.Println("Failed to fetch board games ->", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(boardgames)
}

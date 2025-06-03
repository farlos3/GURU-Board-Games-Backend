package handlers_board

import (
	"log"
	"strconv"

	"guru-game/internal/boardgame/service_board"
	"guru-game/internal/db/repository/boardgame"

	"github.com/gofiber/fiber/v2"
)

// BoardGameHandlers struct สำหรับเก็บ dependencies ของ handlers
type BoardGameHandlers struct {
	BoardGameRepo boardgame.BoardGameRepository
	// ถ้าต้องการเรียก service ที่เชื่อมต่อกับ Python service อาจเพิ่ม field ได้ที่นี่
	// PythonService *service_board.PythonBoardgameService // ตัวอย่าง
}

// NewBoardGameHandlers สร้าง instance ใหม่ของ BoardGameHandlers
func NewBoardGameHandlers(repo boardgame.BoardGameRepository) *BoardGameHandlers {
	return &BoardGameHandlers{
		BoardGameRepo: repo,
	}
}

// HandleGetAllBoardGames handles fetching all board games
func (h *BoardGameHandlers) HandleGetAllBoardGames(c *fiber.Ctx) error {
	// Attempt to get userID from context locals
	userID, ok := c.Locals("userID").(int)
	if !ok {
		userID = 0 // Default to unauthenticated user if no userID is found
	}

	// เรียกใช้ service function โดยตรง
	boardgames, err := service_board.GetAllBoardGames(userID, h.BoardGameRepo) // Pass userID and repo to service
	if err != nil {
		log.Println("Failed to fetch board games ->", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(boardgames)
}

// GetBoardGameByIDHandler handles fetching board game by ID from PostgreSQL
func (h *BoardGameHandlers) GetBoardGameByIDHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Println("Invalid board game ID:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid board game ID"})
	}

	boardgame, err := service_board.GetBoardGameByID(id)
	if err != nil {
		log.Println("Failed to fetch board game ->", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Board game not found"})
	}

	return c.Status(fiber.StatusOK).JSON(boardgame)
}

// GetBoardGameByIDFromESHandler handles fetching board game by ID from Python service
func (h *BoardGameHandlers) GetBoardGameByIDFromESHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Println("Invalid board game ID:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid board game ID"})
	}

	boardgame, err := service_board.GetBoardGameByIDFromES(id)
	if err != nil {
		log.Println("Failed to fetch board game from Python service ->", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Board game not found"})
	}

	return c.Status(fiber.StatusOK).JSON(boardgame)
}

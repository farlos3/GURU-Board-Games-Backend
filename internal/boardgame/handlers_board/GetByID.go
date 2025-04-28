package handlers_board

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"guru-game/internal/boardgame/service_board"
)

// GetBoardGameByIDHandler ดึงข้อมูลบอร์ดเกมตาม ID
func GetBoardGameByIDHandler(c *fiber.Ctx) error {
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

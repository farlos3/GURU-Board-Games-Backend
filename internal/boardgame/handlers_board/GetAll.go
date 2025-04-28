package handlers_board

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"guru-game/internal/boardgame/service_board"
)

// GetAllBoardGamesHandler ดึงข้อมูลบอร์ดเกมทั้งหมด
func GetAllBoardGamesHandler(c *fiber.Ctx) error {
	boardgames, err := service_board.GetAllBoardGames()
	if err != nil {
		log.Println("Failed to fetch board games ->", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(boardgames)
}

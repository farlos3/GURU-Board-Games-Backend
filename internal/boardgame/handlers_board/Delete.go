package handlers_board

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"guru-game/internal/boardgame/service_board"
)

// DeleteBoardGameHandler ลบบอร์ดเกมตาม ID
func DeleteBoardGameHandler(c *fiber.Ctx) error {
	// รับ ID จากพาธ
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Println("Invalid board game ID:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid board game ID"})
	}

	// เรียกฟังก์ชันจาก service เพื่อลบบอร์ดเกม
	err = service_board.DeleteBoardGame(id)
	if err != nil {
		log.Println("Failed to delete board game ->", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Board game deleted successfully"})
}
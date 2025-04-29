package handlers_board

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"guru-game/internal/boardgame/service_board"
	"guru-game/models"
)

// UpdateBoardGameHandler อัปเดตข้อมูลบอร์ดเกม
func UpdateBoardGameHandler(c *fiber.Ctx) error {
	var input models.BoardGame

	// พาร์สข้อมูลที่ได้รับมา
	if err := c.BodyParser(&input); err != nil {
		log.Println("Failed to parse body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// ตรวจสอบข้อมูลที่จำเป็น
	if input.ID == 0 || input.Name == "" || input.Description == "" || input.Players == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Board game ID, name, description, and number of players are required"})
	}

	// เรียกฟังก์ชันจาก service ที่จัดการการอัปเดตข้อมูล boardgame
	updatedBoardGame, err := service_board.UpdateBoardGame(&input)
	if err != nil {
		log.Println("Failed to update board game ->", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(updatedBoardGame)
}
package handlers_board

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"guru-game/internal/boardgame/service_board"
	"guru-game/models"
)

// CreateBoardGameHandler สร้างบอร์ดเกมใหม่
func CreateBoardGameHandler(c *fiber.Ctx) error {
	var input models.BoardGame

	// พาร์สข้อมูลที่ได้รับมา
	if err := c.BodyParser(&input); err != nil {
		log.Println("Failed to parse body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// ตรวจสอบข้อมูลที่จำเป็น
	if input.Name == "" || input.Description == "" || input.Players == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name, description, and number of players are required"})
	}

	// เรียกฟังก์ชันจาก service ที่จัดการการสร้าง boardgame
	boardgame, err := service_board.CreateBoardGame(&input)
	if err != nil {
		log.Println("Failed to create board game ->", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(boardgame)
}

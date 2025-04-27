package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"guru-game/internal/auth/service"
	"guru-game/models"
)

// DeleteUserHandler ลบ user ผ่าน username, email, password
func DeleteUserHandler(c *fiber.Ctx) error {
	var input models.User

	if err := c.BodyParser(&input); err != nil {
		log.Println("Failed to parse body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if input.Username == "" || input.Email == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username, email, and password required"})
	}

	err := service.DeleteUser(&input)
	if err != nil {
		log.Println("Failed to delete user:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User deleted successfully"})
}
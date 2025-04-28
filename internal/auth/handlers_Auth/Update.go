package handlers_Auth

import (
	"log"

    "github.com/gofiber/fiber/v2"
    "guru-game/internal/auth/service_auth"
    "guru-game/models"
)

func UpdateUserHandler(c *fiber.Ctx) error {
	var input models.User

	if err := c.BodyParser(&input); err != nil {
		log.Println("Failed to parse body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// เรียก service ตัวเดียว
	updatedUser, err := service_auth.UpdateUser(&input)
	if err != nil {
		log.Println("Failed to update user:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User updated successfully",
		"user":    updatedUser,
	})
}
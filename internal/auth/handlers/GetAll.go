package handlers

import (
    "github.com/gofiber/fiber/v2"
    "guru-game/internal/auth/service"
)

// GetAllUsers handler
func GetAllUsersHandler(c *fiber.Ctx) error {
	users, err := service.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get users"})
	}
	return c.JSON(users)
}
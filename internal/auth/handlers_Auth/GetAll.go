package handlers_Auth

import (
    "github.com/gofiber/fiber/v2"
    "guru-game/internal/auth/service_auth"
)

// GetAllUsers handler
func GetAllUsersHandler(c *fiber.Ctx) error {
	users, err := service_auth.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get users"})
	}
	return c.JSON(users)
}
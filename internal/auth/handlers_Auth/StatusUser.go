package handlers_Auth

import (
    "github.com/gofiber/fiber/v2"
    "guru-game/models"
)

// GetUser handler
func StatusHandler(c *fiber.Ctx) error {
	user, ok := c.Locals("currentUser").(models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Not logged in"})
	}
	return c.JSON(user)
}
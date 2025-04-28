package handlers_Auth

import (
    "github.com/gofiber/fiber/v2"
    "guru-game/internal/auth/service_auth"
    "guru-game/models"
)

// Login handler
func LoginHandler(c *fiber.Ctx) error {
	input := new(models.User)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	user, err := service_auth.LoginUser(input.Username, input.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	c.Locals("currentUser", *user)
	return c.JSON(fiber.Map{"message": "Login successful", "user": user})
}
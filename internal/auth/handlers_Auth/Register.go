package handlers_Auth

import (
    "github.com/gofiber/fiber/v2"
    "guru-game/internal/auth/service_auth"
    "guru-game/models"
)

// RegisterHandler รับข้อมูลจาก client แล้วส่งไปให้ service จัดการ
func RegisterHandler(c *fiber.Ctx) error {
	newUser := new(models.User)
	if err := c.BodyParser(newUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// เรียกผ่าน service
	user, err := service_auth.RegisterUser(newUser)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Username already exists"})
	}

	return c.JSON(fiber.Map{"message": "User registered", "user": user})
}
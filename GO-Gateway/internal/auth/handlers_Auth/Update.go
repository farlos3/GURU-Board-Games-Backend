package handlers_Auth

import (
	"log"

	"guru-game/internal/auth/service_auth"
	"guru-game/models"

	"github.com/gofiber/fiber/v2"
)

// UpdateUserHandler อัปเดตข้อมูลของผู้ใช้
func UpdateUserHandler(c *fiber.Ctx) error {
	// ดึงข้อมูลจาก JWT ที่เก็บใน context
	user, ok := c.Locals("user").(fiber.Map)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// ดึงข้อมูลจาก request body
	var input models.User
	if err := c.BodyParser(&input); err != nil {
		log.Println("Failed to parse body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// เช็คว่า username หรือ email ถูกต้องหรือไม่
	if input.Username == "" && input.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username or email required"})
	}

	// ตรวจสอบว่า user ที่ต้องการอัปเดตตรงกับ user ที่ได้รับจาก JWT หรือไม่
	if input.Username != "" && input.Username != user["username"].(string) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "You can only update your own information"})
	}

	// เรียกใช้ service เพื่อลงทะเบียนอัปเดตข้อมูลผู้ใช้
	updatedUser, err := service_auth.UpdateUser(&input, user["id"].(int64))
	if err != nil {
		log.Println("Failed to update user ->", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	// ส่งกลับการตอบสนองเมื่ออัปเดตสำเร็จ
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User updated successfully",
		"user":    updatedUser,
	})
}
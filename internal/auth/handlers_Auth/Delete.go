package handlers_Auth

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"guru-game/internal/auth/service_auth"
	"guru-game/models"
)

// DeleteUserHandler ลบ user ผ่าน username, email, password
func DeleteUserHandler(c *fiber.Ctx) error {
	// ดึงข้อมูลจาก JWT
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

	// เช็คว่า username, email, password ถูกต้องหรือไม่
	if input.Username == "" || input.Email == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username, email, and password required"})
	}

	// ตรวจสอบว่าผู้ใช้ที่ทำการลบตรงกับผู้ใช้ที่เป็นเจ้าของ JWT หรือไม่
	if input.Username != user["username"].(string) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "You can only delete your own account"})
	}

	// เรียกใช้ service เพื่อลบ user
	err := service_auth.DeleteUser(&input)
	if err != nil {
		log.Println("Failed to delete user ->", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	// ส่งกลับการตอบสนองเมื่อลบสำเร็จ
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User deleted successfully"})
}
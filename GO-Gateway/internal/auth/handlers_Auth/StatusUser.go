package handlers_Auth

import (
    "github.com/gofiber/fiber/v2"
    "guru-game/models"
)

// StatusHandler ตรวจสอบสถานะของ user ที่ล็อกอิน
func StatusHandler(c *fiber.Ctx) error {
    // ดึงข้อมูล currentUser จาก context
    user, ok := c.Locals("currentUser").(models.User)
    if !ok {
        // ถ้าไม่มีข้อมูลหรือไม่ได้ล็อกอิน
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Not logged in"})
    }

    // ส่งข้อมูล user กลับไปในรูปแบบ JSON
    return c.JSON(fiber.Map{
        "message":  "User is logged in",
        "user_id":  user.ID,
        "username": user.Username,
    })
}

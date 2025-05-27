package handlers_Auth

import (
	"github.com/gofiber/fiber/v2"
	"guru-game/internal/auth/service_auth"
	"guru-game/internal/auth/jwt"
	"guru-game/internal/auth/otp"

	"guru-game/models"
)

// LoginHandler
func LoginHandler(c *fiber.Ctx) error {
    input := new(models.User)
    if err := c.BodyParser(input); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
    }

    // ตรวจสอบ OTP
    otpCode := c.Query("otp") // รับ OTP จาก query parameter หรือ request body
    if otpCode == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "OTP is required"})
    }

    // ตรวจสอบว่า OTP ถูกต้องหรือไม่
    if !otp.VerifyOTP(input.Email, otpCode) {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid OTP"})
    }

    // เรียก service เพื่อทำการ login
    user, err := service_auth.LoginUser(input.Username, input.Password)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
    }

    // สร้าง JWT token สำหรับผู้ใช้ที่ login สำเร็จ
    token, err := jwt.GenerateJWT(user.ID, user.Username)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
    }

    // บันทึกข้อมูล user ลงใน context
    c.Locals("currentUser", *user)

    // ส่ง response กลับไปพร้อมกับ JWT token
    return c.JSON(fiber.Map{
        "message": "Login successful",
        "user":    user,
        "token":   token,
    })
}
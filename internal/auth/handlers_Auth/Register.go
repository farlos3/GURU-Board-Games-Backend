package handlers_Auth

import (
	"github.com/gofiber/fiber/v2"
	"guru-game/internal/auth/otp"
	"guru-game/models"
)

// RegisterHandler
func RegisterHandler(c *fiber.Ctx) error {
	newUser := new(models.User)
	if err := c.BodyParser(newUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// สร้าง OTP
	otpCode, err := otp.GenerateOTP()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate OTP"})
	}

	// ส่ง OTP ไปที่อีเมล
	otp.SendEmail(newUser.Email, otpCode)

	// บันทึก OTP ไปในระบบเพื่อใช้ตรวจสอบในภายหลัง
	otp.SaveOTP(newUser.Email, otpCode)

	// ส่ง OTP กลับไปให้ผู้ใช้เพื่อยืนยัน
	return c.JSON(fiber.Map{
		"message": "OTP sent to your email, please verify",
	})
}
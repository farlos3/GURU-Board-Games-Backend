package handlers_Auth

import (
    "log"

	"github.com/gofiber/fiber/v2"
    
	"guru-game/internal/auth/service_auth"
    "guru-game/internal/auth/otp"
	"guru-game/models"
)

// LoginHandler
func LoginHandler(c *fiber.Ctx) error {
	input := new(models.User)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	log.Println("🔔 [LoginHandler] user submitted:")
	log.Printf("Identifier: %s\n", input.Identifier)

	user, err := service_auth.LoginUser(input.Identifier, input.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// สร้าง OTP ใหม่ทุกครั้งหลัง login ผ่าน
	otpCode, err := otp.GenerateOTP()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate OTP"})
	}

	log.Printf("Generated OTP for %s: %s\n", user.Email, otpCode)

	// บันทึก OTP และ user ชั่วคราว (ในหน่วยความจำ หรือ DB)
	otp.SaveOTP(user.Email, otpCode)
	otp.SaveTempUser(user.Email, *user)

	// ส่ง OTP ทางอีเมล
	otp.SendEmail(user.Email, otpCode)

	// แจ้ง client ว่าต้องยืนยัน OTP
	return c.JSON(fiber.Map{
		"requireOtp": true,
		"email":      user.Email,
		"message":    "OTP sent, please verify it",
	})
}
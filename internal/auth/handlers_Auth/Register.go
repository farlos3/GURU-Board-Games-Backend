package handlers_Auth

import (
	"log"

	"guru-game/models"
	"guru-game/internal/auth/otp"

	"github.com/gofiber/fiber/v2"
	
)

// RegisterHandler
func RegisterHandler(c *fiber.Ctx) error {
	newUser := new(models.User)
	if err := c.BodyParser(newUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	log.Println("🔔 [RegisterHandler] New user submitted:")
	log.Printf("FullName: %s\n", newUser.FullName)
	log.Printf("Username: %s\n", newUser.Username)
	log.Printf("Email: %s\n", newUser.Email)

	// ตรวจสอบว่าอีเมลนี้เคยยืนยันแล้วหรือยัง
	if otp.IsEmailVerified(newUser.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email already verified. Please login or continue registration."})
	}

	// สร้าง OTP, ส่ง OTP, บันทึก OTP ตามเดิม (ถ้าจำเป็น)
	otpCode, err := otp.GenerateOTP()

	log.Printf("OTP: %s\n", otpCode)
	
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate OTP"})
	}
	otp.SendEmail(newUser.Email, otpCode)
	otp.SaveOTP(newUser.Email, otpCode)

	// บันทึกข้อมูลผู้ใช้ชั่วคราว (ในหน่วยความจำ)
	otp.SaveTempUser(newUser.Email, models.User{
		FullName: newUser.FullName,
		Username: newUser.Username,
		Email:    newUser.Email,
		Password: newUser.Password,
	})

	return c.JSON(fiber.Map{
		"message": "OTP sent to your email. Please verify to complete registration.",
	})
}
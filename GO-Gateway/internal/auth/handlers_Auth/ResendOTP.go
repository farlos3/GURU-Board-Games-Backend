package handlers_Auth

import (
	"guru-game/internal/auth/otp"

	"github.com/gofiber/fiber/v2"
)

func ResendOTPHandler(c *fiber.Ctx) error {
	type Request struct {
		Email string `json:"email"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if otp.IsEmailVerified(req.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email already verified"})
	}

	// ✅ รับค่าทั้ง user และ found
	_, found := otp.GetTempUser(req.Email)
	if !found {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No pending registration for this email"})
	}

	otpCode, err := otp.GenerateOTP()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate OTP"})
	}

	otp.SaveOTP(req.Email, otpCode)
	err = otp.SendEmail(req.Email, otpCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send OTP email"})
	}

	return c.JSON(fiber.Map{"message": "New OTP sent to your email"})
}
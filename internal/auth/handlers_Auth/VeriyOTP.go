package handlers_Auth

import (
	"guru-game/internal/auth/otp"
	"guru-game/internal/auth/service_auth"

	"github.com/gofiber/fiber/v2"
)

func VerifyOTPHandler(c *fiber.Ctx) error {
	type Request struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid input",
		})
	}

	if req.Email == "" || req.OTP == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email and OTP are required",
		})
	}

	isValid := otp.VerifyOTP(req.Email, req.OTP)
	if !isValid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid or expired OTP",
		})
	}

	// ดึงข้อมูล user ที่เก็บไว้ชั่วคราว
	user, ok := otp.GetTempUser(req.Email)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "User data not found"})
	}

	// สร้างบัญชีจริง
	createdUser, token, err := service_auth.RegisterUser(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	otp.MarkEmailVerified(req.Email)
	otp.DeleteTempUser(req.Email) // ลบข้อมูลชั่วคราว

	return c.JSON(fiber.Map{
		"message": "OTP verified and user registered successfully.",
		"user":    createdUser,
		"token":   token,
	})
}
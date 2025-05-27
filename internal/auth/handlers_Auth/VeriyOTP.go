package handlers_Auth

import (
	"guru-game/internal/auth/otp"
	"guru-game/internal/auth/service_auth"
	"guru-game/internal/auth/jwt"
	"github.com/gofiber/fiber/v2"
)

func VerifyRegisterOTPHandler(c *fiber.Ctx) error {
	type Request struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid input"})
	}

	if req.Email == "" || req.OTP == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Email and OTP are required"})
	}

	if !otp.VerifyOTP(req.Email, req.OTP) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid or expired OTP"})
	}

	user, ok := otp.GetTempUser(req.Email)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "User data not found"})
	}

	createdUser, token, err := service_auth.RegisterUser(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	otp.MarkEmailVerified(req.Email)
	otp.DeleteTempUser(req.Email)

	return c.JSON(fiber.Map{
		"message": "OTP verified and user registered successfully.",
		"user":    createdUser,
		"token":   token,
	})
}

func VerifyLoginOTPHandler(c *fiber.Ctx) error {
	type Request struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid input"})
	}

	if req.Email == "" || req.OTP == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Email and OTP are required"})
	}

	if !otp.VerifyOTP(req.Email, req.OTP) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid or expired OTP"})
	}

	user, ok := otp.GetTempUser(req.Email)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "User session not found"})
	}

	// สร้าง JWT หลังจาก OTP ถูกต้อง
	token, err := jwt.GenerateJWT(user.ID, user.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	otp.DeleteTempUser(req.Email)

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"user":    user,
		"token":   token,
	})
}

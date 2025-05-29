package handlers_Auth

import (
	"log"

	"guru-game/internal/auth/jwt"
	"guru-game/internal/auth/otp"
	"guru-game/internal/auth/service_auth"

	"github.com/gofiber/fiber/v2"
)

func VerifyRegisterOTPHandler(c *fiber.Ctx) error {
	type Request struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		log.Println("Failed to parse request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if req.Email == "" || req.OTP == "" {
		log.Println("Missing email or OTP in request")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email and OTP are required"})
	}

	// Verify OTP
	if !otp.VerifyOTP(req.Email, req.OTP) {
		log.Printf("Invalid or expired OTP for email: %s", req.Email)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired OTP"})
	}

	// Get temporary user data
	user, ok := otp.GetTempUser(req.Email)
	if !ok {
		log.Printf("No temporary user data found for email: %s", req.Email)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User data not found"})
	}

	// Register user and generate token
	createdUser, token, err := service_auth.RegisterUser(&user)
	if err != nil {
		log.Printf("Failed to register user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Mark email as verified and clean up temporary data
	otp.MarkEmailVerified(req.Email)
	otp.DeleteTempUser(req.Email)

	log.Printf("✅ User registered successfully: %s (ID: %d)", createdUser.Username, createdUser.ID)
	log.Println("🔑 JWT token generated : ", token)

	// Return success response with user data and token
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Registration successful",
		"user": fiber.Map{
			"id":        createdUser.ID,
			"username":  createdUser.Username,
			"email":     createdUser.Email,
			"fullName":  createdUser.FullName,
			"avatarUrl": createdUser.AvatarURL,
			"createdAt": createdUser.CreatedAt,
		},
		"token": token,
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
	log.Println("✅ Login successful")
	log.Println("Token : ", token)

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"user":    user,
		"token":   token,
	})
}

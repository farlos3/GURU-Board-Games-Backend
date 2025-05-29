package handlers_Auth

import (
	"log"

	"guru-game/internal/auth/service_auth"

	"github.com/gofiber/fiber/v2"
)

// GetProfileHandler gets the profile data of the currently logged in user
func GetProfileHandler(c *fiber.Ctx) error {
	// Get user data from JWT token
	user, ok := c.Locals("user").(fiber.Map)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Get user ID from JWT claims
	userID, ok := user["id"].(int64)
	if !ok {
		log.Println("Failed to get user ID from JWT claims")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid user data"})
	}

	// Get username from JWT claims
	username, ok := user["username"].(string)
	if !ok {
		log.Println("Failed to get username from JWT claims")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid user data"})
	}

	// Get user data from database
	userData, err := service_auth.GetUserByID(userID)
	if err != nil {
		log.Println("Failed to get user data:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get user data"})
	}

	// Verify that the user from database matches the JWT claims
	if userData.Username != username {
		log.Println("Username mismatch between JWT and database")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Return user data (excluding sensitive information)
	return c.JSON(fiber.Map{
		"id":        userData.ID,
		"username":  userData.Username,
		"email":     userData.Email,
		"fullName":  userData.FullName,
		"avatarUrl": userData.AvatarURL,
		"createdAt": userData.CreatedAt,
		"updatedAt": userData.UpdatedAt,
	})
}

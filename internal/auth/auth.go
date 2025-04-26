package auth

import (
	"github.com/gofiber/fiber/v2"
	"guru-game/models"
)

// Register handler
func Register(c *fiber.Ctx) error {
	newUser := new(models.User)
	if err := c.BodyParser(newUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	user, err := RegisterUser(newUser) // เรียกใช้ฟังก์ชันใน auth.go (เปลี่ยนจาก service.Register เป็น RegisterUser)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Username already exists"})
	}

	return c.JSON(fiber.Map{"message": "User registered", "user": user})
}

// Login handler
func Login(c *fiber.Ctx) error {
	input := new(models.User)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	user, err := LoginUser(input.Username, input.Password) // เปลี่ยนจาก service.Login เป็น LoginUser
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	c.Locals("currentUser", *user)
	return c.JSON(fiber.Map{"message": "Login successful", "user": user})
}

// GetUser handler
func GetUser(c *fiber.Ctx) error {
	user, ok := c.Locals("currentUser").(models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Not logged in"})
	}
	return c.JSON(user)
}

// GetAllUsers handler
func GetAllUsersHandler(c *fiber.Ctx) error {
	users, err := GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get users"})
	}
	return c.JSON(users)
}
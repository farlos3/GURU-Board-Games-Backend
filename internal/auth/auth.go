package auth

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"guru-game/models"
)

// Register handler
func RegisterHandler(c *fiber.Ctx) error {
	newUser := new(models.User)
	if err := c.BodyParser(newUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	user, err := RegisterUser(newUser)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Username already exists"})
	}

	return c.JSON(fiber.Map{"message": "User registered", "user": user})
}


// Login handler
func LoginHandler(c *fiber.Ctx) error {
	input := new(models.User)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	user, err := LoginUser(input.Username, input.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	c.Locals("currentUser", *user)
	return c.JSON(fiber.Map{"message": "Login successful", "user": user})
}

// GetUser handler
func GetUserHandler(c *fiber.Ctx) error {
	user, ok := c.Locals("currentUser").(models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Not logged in"})
	}
	return c.JSON(user)
}

// UpdateUserHandler รับ request จาก client เพื่อ update user
func UpdateUserHandler(c *fiber.Ctx) error {
	var input models.User

	// Parse ข้อมูลที่ client ส่งมา (เช่น username, email, full_name, avatar_url, password)
	if err := c.BodyParser(&input); err != nil {
		log.Println("❌ Failed to parse body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// ต้องการ username หรือ email อย่างใดอย่างหนึ่ง
	if input.Username == "" && input.Email == "" {
		log.Println("❌ Missing username or email")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username or email required"})
	}

	// หา user ก่อนจาก DB
	var user *models.User
	var err error
	if input.Username != "" {
		user, err = repo.GetByUsername(input.Username)
	} else {
		user, err = repo.GetByEmail(input.Email)
	}

	if err != nil {
		log.Println("❌ User not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// เอา ID ของ user จริงมาใส่ในข้อมูลที่จะอัปเดต
	input.ID = user.ID

	// เรียก Update
	updatedUser, err := UpdateUser(&input)
	if err != nil {
		log.Println("❌ Failed to update user:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "✅ User updated successfully",
		"user":    updatedUser,
	})
}

// GetAllUsers handler
func GetAllUsersHandler(c *fiber.Ctx) error {
	users, err := GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get users"})
	}
	return c.JSON(users)
}
package auth

import (
	"fmt"
	"errors"
	"log"

	"guru-game/models"
	"guru-game/internal/db/repository"
)

var repo db.UserRepository

// Init สำหรับ Inject Repository
func Init(r db.UserRepository) {
	repo = r
}

func RegisterUser(newUser *models.User) (*models.User, error) {
	// ตรวจสอบว่า username ซ้ำไหม
	if user, err := repo.GetByUsername(newUser.Username); err == nil && user != nil {
		log.Printf("Username '%s' already exists.\n", newUser.Username)  // Log if username exists
		return nil, errors.New("username already exists")
	}

	// ตรวจสอบ email ว่ามีค่าไหม ถ้าไม่มีกำหนดเป็นค่า default หรือ error
	if newUser.Email == "" {
		return nil, errors.New("email is required")
	}

	// หากไม่พบ username ซ้ำในฐานข้อมูล, สร้างผู้ใช้ใหม่
	createdUser, err := repo.Create(newUser)
	if err != nil {
		log.Printf("❌ Failed to create user '%s': %v\n", newUser.Username, err)  // Log failure
		return nil, err
	}

	log.Printf("✅ User '%s' created successfully.\n", newUser.Username)  // Log success
	return createdUser, nil
}

func LoginUser(username, password string) (*models.User, error) {
	log.Printf("Attempting to login with username: %s\n", username)

	user, err := repo.GetByCredentials(username, password)
	if err != nil {
		log.Printf("❌ Login failed for username '%s': %v\n", username, err)
		return nil, err
	}

	log.Printf("✅ User '%s' logged in successfully.\n", username) 
	return user, nil
}

// UpdateUser รับ models.User และทำการ update ผ่าน repository
func UpdateUser(updatedUser *models.User) (*models.User, error) {
	log.Printf("🔄 Attempting to update user ID: %d\n", updatedUser.ID)

	// Call repository update
	user, err := repo.Update(updatedUser)
	if err != nil {
		log.Printf("❌ Failed to update user ID %d: %v\n", updatedUser.ID, err)
		return nil, err
	}

	log.Printf("✅ User ID %d updated successfully.\n", updatedUser.ID)
	return user, nil
}

func DeleteUser(userID int64) error {
	log.Printf("Attempting to delete user ID: %d\n", userID)

	err := repo.Delete(userID)
	if err != nil {
		log.Printf("❌ Failed to delete user ID %d: %v\n", userID, err)
		return err
	}

	log.Printf("✅ User ID %d deleted successfully.\n", userID)
	return nil
}


func GetAllUsers() ([]models.User, error) {
	if repo == nil {
		log.Println("User repository is not initialized.")
		return nil, errors.New("user repository is not initialized")
	}

	log.Println("Fetching users from database...")
	users, err := repo.GetAll()
	if err != nil {
		log.Printf("❌ Failed to get users: %v\n", err)
		return nil, fmt.Errorf("failed to get users: %v", err)
	}

	// ถ้าไม่มี user ก็คืน array เปล่า [] ไปเลย ไม่ต้อง return error
	if len(users) == 0 {
		log.Println("No users found.")
		return []models.User{}, nil
	}

	log.Printf("Successfully fetched %d users.\n", len(users))
	return users, nil
}
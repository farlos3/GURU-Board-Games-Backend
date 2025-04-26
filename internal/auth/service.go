package auth

import (
	"fmt"
	"errors"
	"log"  // เพิ่ม log package
	"guru-game/models"
	"guru-game/internal/db/repository"
)

var repo db.UserRepository

// Init สำหรับ Inject Repository
func Init(r db.UserRepository) {
	repo = r
	log.Println("UserRepository initialized successfully.")
}

func RegisterUser(newUser *models.User) (*models.User, error) {
	// ตรวจสอบว่า username ซ้ำไหม
	log.Printf("Attempting to register user with username: %s\n", newUser.Username)

	if user, err := repo.GetByUsername(newUser.Username); err == nil && user != nil {
		log.Printf("Username '%s' already exists.\n", newUser.Username)  // Log if username exists
		return nil, errors.New("username already exists")
	}

	// หากไม่พบ username ซ้ำในฐานข้อมูล, สร้างผู้ใช้ใหม่
	log.Printf("Username '%s' is available, creating user...\n", newUser.Username)
	createdUser, err := repo.Create(newUser)
	if err != nil {
		log.Printf("Failed to create user '%s': %v\n", newUser.Username, err)  // Log failure
		return nil, err
	}

	log.Printf("User '%s' created successfully.\n", newUser.Username)  // Log success
	return createdUser, nil
}

func LoginUser(username, password string) (*models.User, error) {
	log.Printf("Attempting to login with username: %s\n", username)

	user, err := repo.GetByCredentials(username, password)
	if err != nil {
		log.Printf("Login failed for username '%s': %v\n", username, err)
		return nil, err
	}

	log.Printf("User '%s' logged in successfully.\n", username) 
	return user, nil
}

func GetAllUsers() ([]models.User, error) {
	if repo == nil {
		log.Println("User repository is not initialized.")
		return nil, errors.New("user repository is not initialized")
	}

	log.Println("Fetching users from database...")
	users, err := repo.GetAll()
	if err != nil {
		log.Printf("Failed to get users: %v\n", err)
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
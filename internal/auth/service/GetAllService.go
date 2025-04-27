package service

import (
	"errors"
	"log"

	"guru-game/models"
)

func GetAllUsers() ([]models.User, error) {
	if repo == nil {
		log.Println("User repository is not initialized.")
		return nil, errors.New("user repository is not initialized")
	}

	log.Println("Fetching users from database...")
	users, err := repo.GetAll()
	if err != nil {
		log.Printf("Failed to get users: %v\n", err)
		return nil, errors.New("failed to get users: " + err.Error())
	}

	// ถ้าไม่มี user ก็คืน array เปล่า [] ไปเลย ไม่ต้อง return error
	if len(users) == 0 {
		log.Println("No users found.")
		return []models.User{}, nil
	}

	log.Printf("Successfully fetched %d users.\n", len(users))
	return users, nil
}
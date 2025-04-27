package service

import (
	"errors"
	"log"

	"guru-game/models"
)

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
		log.Printf("Failed to create user '%s': %v\n", newUser.Username, err)  // Log failure
		return nil, err
	}

	log.Printf("User '%s' created successfully.\n", newUser.Username)  // Log success
	return createdUser, nil
}
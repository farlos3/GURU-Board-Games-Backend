package service_auth

import (
	"errors"
	"log"

	"guru-game/internal/auth/jwt"
	"guru-game/models"
)

// RegisterUser สมัครผู้ใช้ใหม่และสร้าง JWT token
func RegisterUser(newUser *models.User) (*models.User, string, error) {
	// ตรวจสอบว่า username ซ้ำไหม
	if user, err := repo.GetByUsername(newUser.Username); err == nil && user != nil {
		log.Printf("Username '%s' already exists.\n", newUser.Username)
		return nil, "", errors.New("username already exists")
	}

	// ตรวจสอบ email ว่ามีค่าไหม ถ้าไม่มีกำหนดเป็นค่า default หรือ error
	if newUser.Email == "" {
		return nil, "", errors.New("email is required")
	}

	// หากไม่พบ username ซ้ำในฐานข้อมูล, สร้างผู้ใช้ใหม่
	createdUser, err := repo.Create(newUser)
	if err != nil {
		log.Printf("Failed to create user '%s': %v\n", newUser.Username, err)
		return nil, "", err
	}

	log.Printf("User '%s' created successfully.\n", newUser.Username)

	// สร้าง JWT token หลังจากที่ผู้ใช้ลงทะเบียนสำเร็จ
	token, err := jwt.GenerateJWT(createdUser.ID, createdUser.Username)
	if err != nil {
		log.Printf("Failed to generate JWT for user '%s': %v\n", newUser.Username, err)
		return nil, "", err
	}

	return createdUser, token, nil
}
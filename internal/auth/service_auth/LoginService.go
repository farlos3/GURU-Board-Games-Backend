package service_auth

import (
	"errors"
	"log"

	"guru-game/models"
)

func LoginUser(username, password string) (*models.User, error) {
	log.Printf("🔐 Attempting to login with username: %s\n", username)

	// เช็กว่า username/password ไม่ว่าง
	if username == "" || password == "" {
		log.Println("Username or password is empty")
		return nil, errors.New("username and password must not be empty")
	}

	// เรียก repo หา user
	user, err := repo.GetByCredentials(username, password)
	if err != nil {
		log.Printf("Login failed for username '%s': %v\n", username, err)
		return nil, errors.New("invalid username or password")
	}

	log.Printf("User '%s' logged in successfully.\n", username)
	
	return user, nil
}
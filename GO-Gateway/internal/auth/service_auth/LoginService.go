package service_auth

import (
	"errors"
	"log"

	"guru-game/models"
)

func LoginUser(identifier, password string) (*models.User, error) {
	log.Printf("🔐 Attempting to login with identifier: %s\n", identifier)

	if identifier == "" || password == "" {
		log.Println("Identifier or password is empty")
		return nil, errors.New("identifier and password must not be empty")
	}

	user, err := repo.GetByCredentials(identifier, password)
	if err != nil {
		log.Printf("Login failed for identifier '%s': %v\n", identifier, err)
		return nil, errors.New("invalid email or username or password")
	}

	log.Printf("User '%s' logged in successfully.\n", user.Username)
	return user, nil
}
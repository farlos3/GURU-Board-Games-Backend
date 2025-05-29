package service_auth

import (
	"errors"
	"log"

	"guru-game/models"
)

// GetUserByID retrieves a user by their ID
func GetUserByID(userID int64) (*models.User, error) {
	if repo == nil {
		log.Println("User repository is not initialized")
		return nil, errors.New("user repository is not initialized")
	}

	user, err := repo.GetByID(userID)
	if err != nil {
		log.Printf("Failed to get user with ID %d: %v\n", userID, err)
		return nil, errors.New("failed to get user: " + err.Error())
	}

	return user, nil
}

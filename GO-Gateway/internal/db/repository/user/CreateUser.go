package user

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"guru-game/internal/db/connection"
	"guru-game/models"
)

func (r *PostgresUserRepository) Create(user *models.User) (*models.User, error) {
	// Hash Password with prefix_ + _suffix
	finalPassword := passwordPrefix + user.Password + passwordSuffix
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(finalPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	// updatet current time created_at and updated_at
	currentTime := time.Now()

	query := `INSERT INTO users (username, password, email, full_name, avatar_url, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
		
	err = connection.DB.QueryRow(context.Background(), query, user.Username, string(hashedPassword), user.Email, user.FullName, user.AvatarURL, currentTime, currentTime).Scan(&user.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %v", err)
	}

	user.Password = string(hashedPassword)
	user.CreatedAt = currentTime
	user.UpdatedAt = currentTime

	return user, nil
}
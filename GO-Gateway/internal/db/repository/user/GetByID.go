package user

import (
	"context"
	"fmt"

	"guru-game/internal/db/connection"
	"guru-game/models"
)

// GetByID retrieves a user by their ID
func (r *PostgresUserRepository) GetByID(userID int64) (*models.User, error) {
	query := `SELECT id, username, email, password, full_name, avatar_url, created_at, updated_at FROM users WHERE id = $1`
	row := connection.DB.QueryRow(context.Background(), query, userID)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.FullName, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("user not found by ID: %v", err)
	}
	return &user, nil
}

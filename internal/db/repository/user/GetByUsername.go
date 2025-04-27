package user

import (
	"context"
	"fmt"

	"guru-game/internal/db/connection"
	"guru-game/models"
)

func (r *PostgresUserRepository) GetByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, password, email, full_name, avatar_url, created_at, updated_at FROM users WHERE username = $1`
	row := connection.DB.QueryRow(context.Background(), query, username)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.FullName, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	return &user, nil
}
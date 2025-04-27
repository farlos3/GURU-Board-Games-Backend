package user

import (
	"context"
	"fmt"

	"guru-game/internal/db/connection"
	"guru-game/models"
)

func (r *PostgresUserRepository) GetByEmail(email string) (*models.User, error) {
	query := `SELECT id, username, email, password, full_name, avatar_url, created_at, updated_at FROM users WHERE email = $1`
	row := connection.DB.QueryRow(context.Background(), query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.FullName, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("user not found by email: %v", err)
	}
	return &user, nil
}
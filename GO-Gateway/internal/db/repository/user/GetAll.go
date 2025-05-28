package user

import (
	"context"
	"fmt"

	"guru-game/internal/db/connection"
	"guru-game/models"
)

// Get All User
func (r *PostgresUserRepository) GetAll() ([]models.User, error) {
	query := `SELECT id, username, password, email, full_name, avatar_url, created_at, updated_at FROM users`
	rows, err := connection.DB.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %v", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.FullName, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		users = append(users, user)
	}

	return users, nil
}
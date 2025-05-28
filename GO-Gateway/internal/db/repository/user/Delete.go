package user

import (
	"context"
	"fmt"

	"guru-game/internal/db/connection"
)

// Delete User by ID
func (r *PostgresUserRepository) Delete(userID int64) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := connection.DB.Exec(context.Background(), query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}
	return nil
}
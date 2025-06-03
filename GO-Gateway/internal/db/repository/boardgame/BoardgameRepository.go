package boardgame

import (
	"context"
	"fmt"
	"guru-game/internal/db/connection"
	"guru-game/models"
)

// BoardGameRepository interface for CRUD
type BoardGameRepository interface {
	GetByID(id int) (*models.BoardGame, error)
	GetAll() ([]models.BoardGame, error)
	Delete(id int) error
	GetUserBoardgameState(userID int, boardgameID int) (*models.UserState, error)
}

type PostgresBoardgameRepository struct{}

// GetUserBoardgameState retrieves the user state for a specific board game
func (r *PostgresBoardgameRepository) GetUserBoardgameState(userID int, boardgameID int) (*models.UserState, error) {
	query := `
		SELECT user_id, boardgame_id, liked, favorited, rating, updated_at
		FROM user_states
		WHERE user_id = $1 AND boardgame_id = $2
	`
	row := connection.DB.QueryRow(context.Background(), query, userID, boardgameID)

	var userState models.UserState
	err := row.Scan(
		&userState.UserID,
		&userState.BoardgameID,
		&userState.Liked,
		&userState.Favorited,
		&userState.Rating,
		&userState.UpdatedAt,
	)

	if err != nil {
		// If no rows are found, return nil with no error, or a specific "not found" error
		if err.Error() == "no rows in result set" {
			return nil, nil // Return nil if no state is found
		}
		return nil, fmt.Errorf("failed to fetch user boardgame state: %v", err)
	}

	return &userState, nil
}

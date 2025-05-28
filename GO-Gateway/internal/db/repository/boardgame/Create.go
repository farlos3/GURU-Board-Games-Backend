package boardgame

import (
	"context"
	"fmt"
	"guru-game/internal/db/connection"
	"guru-game/models"
)

// Add Boardgame
func (r *PostgresBoardgameRepository) Create(boardgame *models.BoardGame) (*models.BoardGame, error) {
	query := `INSERT INTO boardgames (name, description, players) VALUES ($1, $2, $3) RETURNING id`
	err := connection.DB.QueryRow(context.Background(), query, boardgame.Name, boardgame.Description, boardgame.Players).Scan(&boardgame.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create board game: %v", err)
	}
	
	return boardgame, nil
}
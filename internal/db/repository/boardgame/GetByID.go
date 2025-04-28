package boardgame

import (
	"context"
	"fmt"
	"guru-game/internal/db/connection"
	"guru-game/models"
)

// ฟังก์ชันค้นหาบอร์ดเกมโดย ID
func (r *PostgresBoardgameRepository) GetByID(id int) (*models.BoardGame, error) {
	query := `SELECT id, name, description, players FROM boardgames WHERE id = $1`
	row := connection.DB.QueryRow(context.Background(), query, id)

	var boardgame models.BoardGame
	err := row.Scan(&boardgame.ID, &boardgame.Name, &boardgame.Description, &boardgame.Players)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("board game with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to fetch board game by ID: %v", err)
	}

	return &boardgame, nil
}
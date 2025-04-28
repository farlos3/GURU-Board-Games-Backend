package boardgame

import (
	"context"
	"fmt"
	"guru-game/internal/db/connection"
	"guru-game/models"
)

// ฟังก์ชันอัปเดตบอร์ดเกมในฐานข้อมูล
func (r *PostgresBoardgameRepository) Update(boardgame *models.BoardGame) (*models.BoardGame, error) {
	query := `UPDATE boardgames SET name = $1, description = $2, players = $3 WHERE id = $4 RETURNING id, name, description, players`
	err := connection.DB.QueryRow(context.Background(), query, boardgame.Name, boardgame.Description, boardgame.Players, boardgame.ID).Scan(&boardgame.ID, &boardgame.Name, &boardgame.Description, &boardgame.Players)
	if err != nil {
		return nil, fmt.Errorf("failed to update board game: %v", err)
	}
	return boardgame, nil
}
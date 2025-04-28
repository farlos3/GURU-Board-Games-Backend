package boardgame

import (
	"context"
	"fmt"
	"guru-game/internal/db/connection"
	"guru-game/models"
)

// ฟังก์ชันดึงข้อมูลบอร์ดเกมทั้งหมดจากฐานข้อมูล
func (r *PostgresBoardgameRepository) GetAll() ([]models.BoardGame, error) {
	query := `SELECT id, name, description, players FROM boardgames`
	rows, err := connection.DB.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch board games: %v", err)
	}
	defer rows.Close()

	var boardgames []models.BoardGame
	for rows.Next() {
		var boardgame models.BoardGame
		err := rows.Scan(&boardgame.ID, &boardgame.Name, &boardgame.Description, &boardgame.Players)
		if err != nil {
			return nil, fmt.Errorf("failed to scan board game: %v", err)
		}
		boardgames = append(boardgames, boardgame)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed iterating rows: %v", err)
	}

	return boardgames, nil
}
package boardgame

import (
	"context"
	"fmt"
	"guru-game/internal/db/connection"
	"guru-game/models"
)

// Get All Boardgame
func (r *PostgresBoardgameRepository) GetAll() ([]models.BoardGame, error) {
	query := `
		SELECT 
			id, title, description, min_players, max_players, play_time_min, play_time_max, 
			categories, rating_avg, rating_count, popularity_score, image_url, created_at, updated_at
		FROM boardgames
	`
	rows, err := connection.DB.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch board games: %v", err)
	}
	defer rows.Close()

	var boardgames []models.BoardGame

	for rows.Next() {
		var bg models.BoardGame
		err := rows.Scan(
			&bg.ID,
			&bg.Title,
			&bg.Description,
			&bg.MinPlayers,
			&bg.MaxPlayers,
			&bg.PlayTimeMin,
			&bg.PlayTimeMax,
			&bg.Categories,
			&bg.RatingAvg,
			&bg.RatingCount,
			&bg.PopularityScore,
			&bg.ImageURL,
			&bg.CreatedAt,
			&bg.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan board game: %v", err)
		}
		boardgames = append(boardgames, bg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed iterating rows: %v", err)
	}

	return boardgames, nil
}

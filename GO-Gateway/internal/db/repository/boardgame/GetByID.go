package boardgame

import (
	"context"
	"fmt"
	"guru-game/internal/db/connection"
	"guru-game/models"
)

// Get Boardgame by ID
func (r *PostgresBoardgameRepository) GetByID(id int) (*models.BoardGame, error) {
	query := `
		SELECT 
			id, title, description, min_players, max_players, 
			play_time_min, play_time_max, categories, rating_avg, rating_count, 
			popularity_score, image_url, created_at, updated_at 
		FROM boardgames 
		WHERE id = $1
	`
	row := connection.DB.QueryRow(context.Background(), query, id)

	var boardgame models.BoardGame
	err := row.Scan(
		&boardgame.ID,
		&boardgame.Title,
		&boardgame.Description,
		&boardgame.MinPlayers,
		&boardgame.MaxPlayers,
		&boardgame.PlayTimeMin,
		&boardgame.PlayTimeMax,
		&boardgame.Categories,
		&boardgame.RatingAvg,
		&boardgame.RatingCount,
		&boardgame.PopularityScore,
		&boardgame.ImageURL,
		&boardgame.CreatedAt,
		&boardgame.UpdatedAt,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("board game with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to fetch board game by ID: %v", err)
	}

	return &boardgame, nil
}
package game_rules

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"guru-game/models"
)

// GameRuleRepository defines the interface for game rules database operations
type GameRuleRepository interface {
	GetRulesByBoardgameID(ctx context.Context, boardgameID int) ([]models.GameRule, error)
}

// PostgresGameRuleRepository handles database operations for GameRule using pgxpool
type PostgresGameRuleRepository struct {
	DB *pgxpool.Pool
}

// NewPostgresGameRuleRepository creates a new PostgresGameRuleRepository
func NewPostgresGameRuleRepository(db *pgxpool.Pool) *PostgresGameRuleRepository {
	return &PostgresGameRuleRepository{DB: db}
}

// GetRulesByBoardgameID fetches game rules for a given boardgame ID from PostgreSQL
func (r *PostgresGameRuleRepository) GetRulesByBoardgameID(ctx context.Context, boardgameID int) ([]models.GameRule, error) {
	query := `
		SELECT id, boardgame_id, title, steps
		FROM game_rules
		WHERE boardgame_id = $1
	`

	rows, err := r.DB.Query(ctx, query, boardgameID)
	if err != nil {
		log.Printf("Error fetching game rules for boardgame %d: %v", boardgameID, err)
		return nil, fmt.Errorf("failed to fetch game rules: %w", err)
	}
	defer rows.Close()

	var gameRules []models.GameRule
	for rows.Next() {
		var rule models.GameRule
		err := rows.Scan(&rule.ID, &rule.BoardgameID, &rule.Title, &rule.Steps)
		if err != nil {
			log.Printf("Error scanning game rule row for boardgame %d: %v", boardgameID, err)
			return nil, fmt.Errorf("failed to scan game rule row: %w", err)
		}
		gameRules = append(gameRules, rule)
	}

	// Check for errors after iterating over rows
	if err := rows.Err(); err != nil {
		log.Printf("Error after iterating over game rule rows for boardgame %d: %v", boardgameID, err)
		return nil, fmt.Errorf("error after fetching game rules: %w", err)
	}

	log.Printf("Successfully fetched %d game rules for boardgame %d", len(gameRules), boardgameID)

	return gameRules, nil
}

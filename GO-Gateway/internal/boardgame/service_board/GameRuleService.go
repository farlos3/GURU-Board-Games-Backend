package service_board

import (
	"context"
	"fmt"
	"log"

	"guru-game/internal/db/repository/game_rules"
	"guru-game/models"
)

// GameRuleService provides access to game rules data
type GameRuleService struct {
	repo game_rules.GameRuleRepository
}

// NewGameRuleService creates a new GameRuleService
func NewGameRuleService(repo game_rules.GameRuleRepository) *GameRuleService {
	return &GameRuleService{repo: repo}
}

// GetGameRulesByBoardgameID fetches game rules for a specific boardgame ID
func (s *GameRuleService) GetGameRulesByBoardgameID(ctx context.Context, boardgameID int) ([]models.GameRule, error) {
	gameRules, err := s.repo.GetRulesByBoardgameID(ctx, boardgameID)
	if err != nil {
		log.Printf("Error from repository when fetching game rules for boardgame %d: %v", boardgameID, err)
		return nil, fmt.Errorf("failed to fetch game rules from repository: %w", err)
	}

	// You might add business logic here if needed

	return gameRules, nil
}

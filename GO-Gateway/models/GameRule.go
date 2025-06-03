package models

// GameRule represents the structure of the game_rules table
type GameRule struct {
	ID          int         `json:"id"`
	BoardgameID int         `json:"boardgame_id"`
	Title       string      `json:"title"`
	Steps       interface{} `json:"steps"` // Using interface{} for jsonb/json type
}

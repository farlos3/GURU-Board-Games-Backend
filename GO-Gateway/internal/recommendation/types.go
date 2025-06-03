package recommendation

import "time"

type UserAction struct {
	UserID      string    `json:"user_id"`
	BoardgameID string    `json:"boardgame_id"`
	ActionType  string    `json:"action_type"`
	ActionValue float64   `json:"action_value"`
	ActionTime  time.Time `json:"action_time"`
}

type Boardgame struct {
	ID              int     `json:"id"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	MinPlayers      int     `json:"min_players"`
	MaxPlayers      int     `json:"max_players"`
	PlayTimeMin     int     `json:"play_time_min"`
	PlayTimeMax     int     `json:"play_time_max"`
	Categories      string  `json:"categories"`
	RatingAvg       float64 `json:"rating_avg"`
	RatingCount     int     `json:"rating_count"`
	PopularityScore float64 `json:"popularity_score"`
	ImageURL        string  `json:"image_url"`
}

package models

import "time"

type BoardGame struct {
	ID              int       `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	MinPlayers      int       `json:"min_players"`
	MaxPlayers      int       `json:"max_players"`
	PlayTimeMin     int       `json:"play_time_min"`
	PlayTimeMax     int       `json:"play_time_max"`
	Categories      string    `json:"categories"` // หรือ slice ถ้าแปลง JSON ใน DB
	RatingAvg       float64   `json:"rating_avg"`
	RatingCount     int       `json:"rating_count"`
	PopularityScore float64   `json:"popularity_score"`
	ImageURL        string    `json:"image_url"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// User-specific state
	LikedByCurrentUser     bool    `json:"likedByCurrentUser,omitempty"`
	FavoritedByCurrentUser bool    `json:"favoritedByCurrentUser,omitempty"`
	CurrentUserRating      float64 `json:"currentUserRating,omitempty"`
}

// ActivityData represents the nested data structure within the request body
type ActivityData struct {
	GameID      int     `json:"gameID"`
	IsLiked     bool    `json:"isLiked"`
	IsViewed    bool    `json:"isViewed"`
	IsPlayed    bool    `json:"isPlayed"`
	IsFavorite  bool    `json:"isFavorite"`
	RatingValue float64 `json:"ratingValue"`
	SearchQuery string  `json:"searchQuery"`
}

// GameSearchQuery represents the expected structure of the incoming query parameters
type GameSearchQuery struct {
	SearchQuery string   `query:"SearchQuery"`
	Categories  []string `query:"categories"`
	PlayerCount int      `query:"playerCount"`
	PlayTime    int      `query:"playTime"`
	Limit       int      `query:"limit"`
	Page        int      `query:"page"`
}

// UserState represents a row in the user_states table
type UserState struct {
	UserID      int       `json:"user_id"`
	BoardgameID int       `json:"boardgame_id"`
	Liked       bool      `json:"liked"`
	Favorited   bool      `json:"favorited"`
	Rating      float64   `json:"rating"`
	UpdatedAt   time.Time `json:"updated_at"`
}

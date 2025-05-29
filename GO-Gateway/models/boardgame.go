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
	Categories      string    `json:"categories"`       // หรือ slice ถ้าแปลง JSON ใน DB
	RatingAvg       float64   `json:"rating_avg"`
	RatingCount     int       `json:"rating_count"`
	PopularityScore float64   `json:"popularity_score"`
	ImageURL        string    `json:"image_url"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
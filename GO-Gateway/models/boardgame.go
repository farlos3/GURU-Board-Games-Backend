package models

type BoardGame struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Players     int    `json:"players"`
}
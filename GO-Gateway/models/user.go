package models

import "time"

// เก็บข้อมูลใน Database
type User struct {
	ID         int64     `json:"id"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	Email      string    `json:"email"`
	FullName   string    `json:"fullName"`
	AvatarURL  string    `json:"avatar_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Identifier string    `json:"identifier"`

	RequireOtp bool `json:"requireOtp"`
}

// ActivityLog represents the expected structure of the incoming request body
type ActivityLog struct {
	Type      string       `json:"type"`
	Data      ActivityData `json:"data"`
	UserID    int64        `json:"userID"`    // Changed to int64 to match incoming data
	SessionID string       `json:"sessionID"` // Assuming sessionID is a string, adjust if needed
	Timestamp string       `json:"timestamp"` // Assuming timestamp is a string, adjust if needed
}
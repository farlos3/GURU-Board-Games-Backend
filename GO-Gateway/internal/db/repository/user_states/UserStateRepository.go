package user_states

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// UserState represents the user_states table schema
type UserState struct {
	UserID      int       `json:"user_id"`
	BoardgameID int       `json:"boardgame_id"`
	Liked       bool      `json:"liked"`
	Favorited   bool      `json:"favorited"`
	Rating      float64   `json:"rating"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserStateRepository defines the interface for user state database operations
type UserStateRepository interface {
	SaveOrUpdate(ctx context.Context, userState *UserState) error
	GetFavoritedByUserID(ctx context.Context, userID int) ([]UserState, error)
	GetAllByUserID(ctx context.Context, userID int) ([]UserState, error)
}

// PostgresUserStateRepository handles database operations for UserState using pgxpool
type PostgresUserStateRepository struct {
	DB *pgxpool.Pool
}

// NewPostgresUserStateRepository creates a new PostgresUserStateRepository
func NewPostgresUserStateRepository(db *pgxpool.Pool) *PostgresUserStateRepository {
	return &PostgresUserStateRepository{DB: db}
}

// SaveOrUpdate saves or updates a user state in PostgreSQL
func (r *PostgresUserStateRepository) SaveOrUpdate(ctx context.Context, userState *UserState) error {
	// Use INSERT ... ON CONFLICT (user_id, boardgame_id) DO UPDATE to upsert
	query := `
		INSERT INTO user_states (user_id, boardgame_id, liked, favorited, rating, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, boardgame_id) DO UPDATE
		SET liked = EXCLUDED.liked,
		    favorited = EXCLUDED.favorited,
		    rating = EXCLUDED.rating,
		    updated_at = EXCLUDED.updated_at
	`

	_, err := r.DB.Exec(ctx, query,
		userState.UserID,
		userState.BoardgameID,
		userState.Liked,
		userState.Favorited,
		userState.Rating,
		time.Now(),
	)

	if err != nil {
		log.Printf("Error saving or updating user state: %v", err)
		return fmt.Errorf("failed to save or update user state: %w", err)
	}

	log.Printf("User state saved/updated successfully for user_id: %d, boardgame_id: %d", userState.UserID, userState.BoardgameID)

	return nil
}

// GetFavoritedByUserID fetches favorited user states for a given user ID from PostgreSQL
func (r *PostgresUserStateRepository) GetFavoritedByUserID(ctx context.Context, userID int) ([]UserState, error) {
	query := `
		SELECT user_id, boardgame_id, liked, favorited, rating, updated_at
		FROM user_states
		WHERE user_id = $1 AND favorited = TRUE
	`

	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		log.Printf("Error fetching favorited user states for user %d: %v", userID, err)
		return nil, fmt.Errorf("failed to fetch favorited user states: %w", err)
	}
	defer rows.Close()

	var favoritedStates []UserState
	for rows.Next() {
		var state UserState
		err := rows.Scan(&state.UserID, &state.BoardgameID, &state.Liked, &state.Favorited, &state.Rating, &state.UpdatedAt)
		if err != nil {
			log.Printf("Error scanning favorited user state row: %v", err)
			return nil, fmt.Errorf("failed to scan favorited user state row: %w", err)
		}
		favoritedStates = append(favoritedStates, state)
	}

	// Check for errors after iterating over rows
	if err := rows.Err(); err != nil {
		log.Printf("Error after iterating over favorited user state rows: %v", err)
		return nil, fmt.Errorf("error after fetching favorited user states: %w", err)
	}

	log.Printf("Successfully fetched %d favorited user states for user %d", len(favoritedStates), userID)

	return favoritedStates, nil
}

// GetAllByUserID fetches all user states for a given user ID from PostgreSQL
func (r *PostgresUserStateRepository) GetAllByUserID(ctx context.Context, userID int) ([]UserState, error) {
	query := `
		SELECT user_id, boardgame_id, liked, favorited, rating, updated_at
		FROM user_states
		WHERE user_id = $1
	`

	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		log.Printf("Error fetching user states for user %d: %v", userID, err)
		return nil, fmt.Errorf("failed to fetch user states: %w", err)
	}
	defer rows.Close()

	var userStates []UserState
	for rows.Next() {
		var state UserState
		err := rows.Scan(&state.UserID, &state.BoardgameID, &state.Liked, &state.Favorited, &state.Rating, &state.UpdatedAt)
		if err != nil {
			log.Printf("Error scanning user state row: %v", err)
			return nil, fmt.Errorf("failed to scan user state row: %w", err)
		}
		userStates = append(userStates, state)
	}

	// Check for errors after iterating over rows
	if err := rows.Err(); err != nil {
		log.Printf("Error after iterating over user state rows: %v", err)
		return nil, fmt.Errorf("error after fetching user states: %w", err)
	}

	log.Printf("Successfully fetched %d user states for user %d", len(userStates), userID)

	return userStates, nil
}

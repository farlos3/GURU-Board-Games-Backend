package handlers

import (
	"context"
	"guru-game/internal/db/connection"
	"guru-game/internal/recommendation"
	"guru-game/models"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Handler for user activity
type UserActivityHandler struct {
	recommendationClient recommendation.RecommendationClient
}

// NewUserActivityHandler creates a new handler instance
func NewUserActivityHandler(client recommendation.RecommendationClient) *UserActivityHandler {
	return &UserActivityHandler{
		recommendationClient: client,
	}
}

// HandleUserActivity receives and processes user activity logs
func (h *UserActivityHandler) HandleUserActivity(c *fiber.Ctx) error {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("\n\n[%s] ===== USER ACTIVITY LOG RECEIVED =====\n", timestamp)

	// Log raw request body
	rawBody := c.Body()
	log.Printf("[%s] Raw Request Body: %s", timestamp, string(rawBody))

	var activityLog models.ActivityLog

	// Parse the request body into the ActivityLog struct
	if err := c.BodyParser(&activityLog); err != nil {
		log.Printf("[%s] Error parsing activity log request body: %v", timestamp, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	// Log the parsed data in a structured format - logging removed to simplify output for now

	// --- Add logic to process activity and update user_states table ---
	ctx := context.Background()

	switch activityLog.Type {
	case "LIKE_GAME":
		log.Printf("[%s] Processing LIKE_GAME for UserID: %d, GameID: %d, IsLiked: %v", timestamp, activityLog.UserID, activityLog.Data.GameID, activityLog.Data.IsLiked)
		// UPSERT into user_states table
		query := `INSERT INTO user_states (user_id, boardgame_id, liked, updated_at)
		          VALUES ($1, $2, $3, NOW())
		          ON CONFLICT (user_id, boardgame_id) DO UPDATE SET
		          liked = EXCLUDED.liked, updated_at = NOW()`
		_, err := connection.DB.Exec(ctx, query, activityLog.UserID, activityLog.Data.GameID, activityLog.Data.IsLiked)
		if err != nil {
			log.Printf("[%s] Error updating user_states for LIKE_GAME: %v", timestamp, err)
			// Decide how to handle error - return error to frontend?
		}

	case "FAVORITE_GAME":
		log.Printf("[%s] Processing FAVORITE_GAME for UserID: %d, GameID: %d, IsFavorite: %v", timestamp, activityLog.UserID, activityLog.Data.GameID, activityLog.Data.IsFavorite)
		// UPSERT into user_states table
		query := `INSERT INTO user_states (user_id, boardgame_id, favorited, updated_at)
		          VALUES ($1, $2, $3, NOW())
		          ON CONFLICT (user_id, boardgame_id) DO UPDATE SET
		          favorited = EXCLUDED.favorited, updated_at = NOW()`
		_, err := connection.DB.Exec(ctx, query, activityLog.UserID, activityLog.Data.GameID, activityLog.Data.IsFavorite)
		if err != nil {
			log.Printf("[%s] Error updating user_states for FAVORITE_GAME: %v", timestamp, err)
			// Decide how to handle error
		}

	case "RATE_GAME":
		log.Printf("[%s] Processing RATE_GAME for UserID: %d, GameID: %d, RatingValue: %v", timestamp, activityLog.UserID, activityLog.Data.GameID, activityLog.Data.RatingValue)
		// UPSERT into user_states table
		query := `INSERT INTO user_states (user_id, boardgame_id, rating, updated_at)
		          VALUES ($1, $2, $3, NOW())
		          ON CONFLICT (user_id, boardgame_id) DO UPDATE SET
		          rating = EXCLUDED.rating, updated_at = NOW()`
		_, err := connection.DB.Exec(ctx, query, activityLog.UserID, activityLog.Data.GameID, activityLog.Data.RatingValue)
		if err != nil {
			log.Printf("[%s] Error updating user_states for RATE_GAME: %v", timestamp, err)
			// Decide how to handle error
		}

	// TODO: Add cases for other activity types like VIEW_GAME, PLAY_GAME, etc.

	default:
		log.Printf("[%s] Unhandled activity type: %s", timestamp, activityLog.Type)
		// Optionally handle other types or return a different response
	}

	// Note: Sending to Recommendation Service logic is temporarily commented out
	// if activityLog.Type in types to be sent to recommendation service...
	// ... send to h.recommendationClient ...

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Activity log received and processed successfully",
	})
}
package handlers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GameStateUpdateData represents the nested data structure for updates
type GameStateUpdateData map[string]interface{}

// GameStateUpdate represents the expected structure of the incoming request body
type GameStateUpdate struct {
	UserID     string              `json:"userID"`
	GameID     string              `json:"gameID"`
	UpdateData GameStateUpdateData `json:"updateData"`
}

// HandleGameStateUpdate receives and processes game state updates
func HandleGameStateUpdate(c *fiber.Ctx) error {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("\n\n[%s] ===== GAME STATE UPDATE RECEIVED =====\n", timestamp)

	// Log request headers
	log.Printf("[%s] Request Headers:", timestamp)
	for k, v := range c.GetReqHeaders() {
		log.Printf("[%s]   %s: %s", timestamp, k, v)
	}

	// Log raw request body
	rawBody := c.Body()
	log.Printf("[%s] Raw Request Body: %s", timestamp, string(rawBody))

	var update GameStateUpdate

	// Parse the request body
	if err := c.BodyParser(&update); err != nil {
		log.Printf("[%s] Error parsing game state update request body: %v", timestamp, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	// Log the parsed data
	log.Printf("[%s] Parsed Game State Update:", timestamp)
	log.Printf("[%s]   UserID: %s", timestamp, update.UserID)
	log.Printf("[%s]   GameID: %s", timestamp, update.GameID)
	log.Printf("[%s]   UpdateData:", timestamp)
	for k, v := range update.UpdateData {
		log.Printf("[%s]     %s: %v", timestamp, k, v)
	}

	// TODO: Implement further processing (e.g., update database, interact with other services)

	log.Printf("[%s] ===== GAME STATE UPDATE PROCESSED =====\n", timestamp)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Game state update received successfully",
	})
}

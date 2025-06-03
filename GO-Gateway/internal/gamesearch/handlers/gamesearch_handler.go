package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"guru-game/models"
)

// HandleGameSearch receives and processes game search queries
func HandleGameSearch(c *fiber.Ctx) error {
	log.Println("Received game search query")

	var query models.GameSearchQuery

	// Parse the query parameters into the GameSearchQuery struct
	if err := c.QueryParser(&query); err != nil {
		log.Printf("Error parsing game search query parameters: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse query parameters",
		})
	}

	// Log the received data for now
	log.Printf("Received Game Search Query: %+v", query)

	// TODO: Implement further processing (e.g., query ElasticSearch, return results)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Game search query received successfully",
		"query":   query,
	})
}

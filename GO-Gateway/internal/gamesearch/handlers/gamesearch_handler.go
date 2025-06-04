package handlers

import (
	"log"

	"encoding/json"
	"fmt"
	"guru-game/models"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// GameSearchHandlers holds the necessary dependencies for game search handlers
type GameSearchHandlers struct {
	PythonServiceURL string
}

// NewGameSearchHandlers creates a new GameSearchHandlers instance
func NewGameSearchHandlers(pythonServiceURL string) *GameSearchHandlers {
	return &GameSearchHandlers{PythonServiceURL: pythonServiceURL}
}

// HandleGameSearch receives and processes game search queries by forwarding to Python service
func (h *GameSearchHandlers) HandleGameSearch(c *fiber.Ctx) error {
	log.Println("Received game search query")

	var query models.GameSearchQuery

	// Parse the query parameters into the GameSearchQuery struct
	if err := c.QueryParser(&query); err != nil {
		log.Printf("Error parsing game search query parameters: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse query parameters",
		})
	}

	log.Printf("Received Game Search Query: %+v", query)

	// Construct the URL for the Python service search endpoint
	pythonSearchURL := fmt.Sprintf("%s/api/search?searchQuery=%s&playerCount=%d&playTime=%d",
		h.PythonServiceURL, query.SearchQuery, query.PlayerCount, query.PlayTime)

	// Add categories if selected
	if len(query.Categories) > 0 {
		pythonSearchURL = fmt.Sprintf("%s&categories=%s", pythonSearchURL, strings.Join(query.Categories, ","))
	}

	// Add limit and page for pagination
	if query.Limit > 0 {
		pythonSearchURL = fmt.Sprintf("%s&limit=%d", pythonSearchURL, query.Limit)
	}
	if query.Page > 0 {
		pythonSearchURL = fmt.Sprintf("%s&page=%d", pythonSearchURL, query.Page)
	}

	log.Printf("Forwarding search request to Python service: %s", pythonSearchURL)

	// Make the HTTP GET request to the Python service
	resp, err := http.Get(pythonSearchURL)
	if err != nil {
		log.Printf("Error forwarding search request to Python service: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to connect to search service",
		})
	}
	defer resp.Body.Close()

	// Read the response body
	var searchResults []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&searchResults); err != nil {
		log.Printf("Error decoding response from Python service: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process search results",
		})
	}

	// Return the results from the Python service to the frontend
	return c.Status(resp.StatusCode).JSON(searchResults)
}

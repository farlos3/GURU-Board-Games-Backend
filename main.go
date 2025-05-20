package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"

	"guru-game/internal/auth/service_auth"
	"guru-game/internal/boardgame/service_board"

	"guru-game/internal/db/connection"
	"guru-game/internal/db/repository/boardgame"
	"guru-game/internal/db/repository/user"
	"guru-game/routes"
)

func main() {
	// Start fiber
	app := fiber.New()

	// Connect DB
	connection.ConnectDB()
	service_auth.Init(&user.PostgresUserRepository{})
	service_board.Init(&boardgame.PostgresBoardgameRepository{})

	// Set up routes
	routes.SetupRoutes(app)

	// Read port from environment or default to 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 Server is running on http://localhost:%s\n", port)
	log.Fatal(app.Listen(":" + port))
}
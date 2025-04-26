package routes

import (
	"github.com/gofiber/fiber/v2"
	"guru-game/internal/auth"
	"guru-game/internal/boardgame"
)

// SetupRoutes initializes all API routes
func SetupRoutes(app *fiber.App) {
	// Auth routes
	api := app.Group("/auth")
	api.Post("/register", auth.RegisterHandler)
	api.Post("/login", auth.LoginHandler)
	api.Get("/user", auth.GetUserHandler)
	api.Get("/users", auth.GetAllUsersHandler)
	api.Put("/user/update", auth.UpdateUserHandler)

	// Boardgame routes
	bg := app.Group("/boardgames")
	bg.Post("/add", boardgame.CreateBoardGame)        
	bg.Get("/", boardgame.GetAllBoardGames)        
	bg.Get("/:name", boardgame.GetBoardGameByName) 
}
package routes

import (
	"github.com/gofiber/fiber/v2"
	"guru-game/internal/auth/handlers"
	// "guru-game/internal/boardgame"
)

// SetupRoutes initializes all API routes
func SetupRoutes(app *fiber.App) {
	// Auth routes
	api := app.Group("/auth")
	api.Post("/register", handlers.RegisterHandler)
	api.Post("/login", handlers.LoginHandler)
	api.Get("/status", handlers.StatusHandler)
	api.Get("/users", handlers.GetAllUsersHandler)
	api.Put("/user/update", handlers.UpdateUserHandler)
	api.Delete("/user/delete", handlers.DeleteUserHandler)

	// Boardgame routes
	// bg := app.Group("/boardgames")
	// bg.Post("/add", boardgame.CreateBoardGame)        
	// bg.Get("/", boardgame.GetAllBoardGames)        
	// bg.Get("/:name", boardgame.GetBoardGameByName) 
}
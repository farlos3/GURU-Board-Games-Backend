package routes

import (
	"github.com/gofiber/fiber/v2"
	"guru-game/internal/auth"
	"guru-game/internal/boardgame" // เพิ่มตรงนี้ด้วย
)

// SetupRoutes initializes all API routes
func SetupRoutes(app *fiber.App) {
	// Auth routes
	api := app.Group("/auth")
	api.Post("/register", auth.Register)
	api.Post("/login", auth.Login)
	api.Get("/user", auth.GetUser)
	api.Get("/users", auth.GetAllUsersHandler)

	// Boardgame routes
	bg := app.Group("/boardgames")
	bg.Post("/add", boardgame.CreateBoardGame)        
	bg.Get("/", boardgame.GetAllBoardGames)        
	bg.Get("/:name", boardgame.GetBoardGameByName) 
}
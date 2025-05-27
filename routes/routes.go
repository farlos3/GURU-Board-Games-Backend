package routes

import (
	"github.com/gofiber/fiber/v2"
	"guru-game/internal/auth/handlers_Auth"
	"guru-game/internal/boardgame/handlers_board"
	"guru-game/internal/auth/jwt"
)

func SetupRoutes(app *fiber.App) {
	// Auth routes
	api := app.Group("/auth")
	api.Post("/register", handlers_Auth.RegisterHandler)
	api.Post("/login", handlers_Auth.LoginHandler)
	app.Post("/auth/verify-otp", handlers_Auth.VerifyOTPHandler)

	api.Get("/status", jwt.JWTMiddleware, handlers_Auth.StatusHandler)
	api.Get("/users", handlers_Auth.GetAllUsersHandler)
	api.Put("/user/update", jwt.JWTMiddleware, handlers_Auth.UpdateUserHandler)
	api.Delete("/user/delete", jwt.JWTMiddleware, handlers_Auth.DeleteUserHandler)

	// Boardgame routes
	bg := app.Group("/boardgames")
	bg.Post("/add", handlers_board.CreateBoardGameHandler)       
	bg.Get("/", handlers_board.GetAllBoardGamesHandler)          
	bg.Get("/:id", handlers_board.GetBoardGameByIDHandler)       
	bg.Put("/:id", handlers_board.UpdateBoardGameHandler)
	bg.Delete("/:id", handlers_board.DeleteBoardGameHandler)
}
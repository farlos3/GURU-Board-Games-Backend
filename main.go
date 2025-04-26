package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"guru-game/routes"
	"guru-game/internal/db"
	"guru-game/internal/auth"
	"guru-game/internal/boardgame"
)

func main() {
	app := fiber.New()

	// Init mock database
	db.ConnectMock()      // สำหรับ User
	db.ConnectMockGame()  // สำหรับ BoardGame

	// Init Repositories
	auth.Init(db.MockUserRepository{})         // auth ใช้ repo ของ User
	boardgame.Init(db.MockBoardgameRepository{}) // boardgame ใช้ repo ของ Boardgame

	// Setup all routes (auth + boardgame)
	routes.SetupRoutes(app)

	log.Println("🚀 Server is running on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}
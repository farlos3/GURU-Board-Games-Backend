package main

import (
	"log"
	"os"

	"guru-game/internal/auth/service_auth"
	"guru-game/internal/boardgame/service_board"

	"guru-game/internal/db/connection"
	"guru-game/internal/db/repository/boardgame"
	"guru-game/internal/db/repository/user"
	"guru-game/routes"
	
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	log.Println("🚀 Starting server...")
	
	app := fiber.New()

	// ตั้งค่า CORS middleware เพื่ออนุญาต frontend จาก localhost:3000
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	log.Println("✅ CORS middleware configured")

	// เชื่อมต่อ DB
	connection.ConnectDB()
	service_auth.Init(&user.PostgresUserRepository{})
	service_board.Init(&boardgame.PostgresBoardgameRepository{})
	log.Println("✅ Services initialized")

	// ตั้งค่า routes
	log.Println("🔧 Setting up routes...")
	routes.SetupRoutes(app)
	log.Println("✅ Routes configured")

	port := os.Getenv("GO_PORT")
	if port == "" {
		port = "5000"
	}

	log.Printf("🚀 Server is running on http://localhost:%s\n", port)
	log.Fatal(app.Listen(":" + port))
}
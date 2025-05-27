package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"guru-game/internal/auth/service_auth"
	"guru-game/internal/boardgame/service_board"

	"guru-game/internal/db/connection"
	"guru-game/internal/db/repository/boardgame"
	"guru-game/internal/db/repository/user"
	"guru-game/routes"
)

func main() {
	app := fiber.New()

	// ตั้งค่า CORS middleware เพื่ออนุญาต frontend จาก localhost:3000
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000", // เปลี่ยนเป็นโดเมน frontend ของคุณ
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Connect DB
	connection.ConnectDB()
	service_auth.Init(&user.PostgresUserRepository{})
	service_board.Init(&boardgame.PostgresBoardgameRepository{})

	// Set up routes
	routes.SetupRoutes(app)

	port := os.Getenv("GO_PORT")
	if port == "" {
		port = "5000"
	}

	log.Printf("🚀 Server is running on http://localhost:%s\n", port)
	log.Fatal(app.Listen(":" + port))
}

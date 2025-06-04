package main

import (
	"log"
	"os"

	"guru-game/internal/auth/service_auth"
	"guru-game/internal/boardgame/service_board"

	"guru-game/internal/db/connection"
	"guru-game/internal/db/repository/boardgame"
	"guru-game/internal/db/repository/game_rules"
	"guru-game/internal/db/repository/user"
	"guru-game/internal/db/repository/user_states"
	"guru-game/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	// Ensure the gamesearch handlers package is imported
	gamesearchhandlers "guru-game/internal/gamesearch/handlers"
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

	// Initialize repositories
	userStateRepo := user_states.NewPostgresUserStateRepository(connection.DB)
	// Initialize boardGameRepo correctly as an empty struct
	boardGameRepo := &boardgame.PostgresBoardgameRepository{}
	gameRuleRepo := game_rules.NewPostgresGameRuleRepository(connection.DB)
	log.Println("✅ Repositories initialized")

	// Initialize services
	// Provide the boardGameRepo to the boardgame service
	service_board.Init(boardGameRepo)
	gameRuleService := service_board.NewGameRuleService(gameRuleRepo)
	log.Println("✅ Services initialized")

	// Initialize Game Search Handlers
	pythonServiceURL := os.Getenv("PYTHON_SERVICE_URL")
	if pythonServiceURL == "" {
		pythonServiceURL = "http://localhost:50051" // default URL
	}
	gameSearchHandlers := gamesearchhandlers.NewGameSearchHandlers(pythonServiceURL)

	log.Println("🔧 Setting up routes...")
	// Pass the concrete boardGameRepo which satisfies the interface
	routes.SetupRoutes(app, userStateRepo, boardGameRepo, gameRuleService, gameSearchHandlers)
	log.Println("✅ Routes configured")

	port := os.Getenv("GO_PORT")
	if port == "" {
		port = "5000"
	}

	log.Printf("🚀 Server is running on http://localhost:%s\n", port)
	log.Fatal(app.Listen(":" + port))
}

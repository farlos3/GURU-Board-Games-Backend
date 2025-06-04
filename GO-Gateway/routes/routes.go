package routes

import (
	"guru-game/internal/auth/handlers_Auth"
	"guru-game/internal/auth/jwt"
	"guru-game/internal/boardgame/handlers_board"
	"guru-game/internal/boardgame/service_board"
	"guru-game/internal/db/repository/boardgame"
	"guru-game/internal/db/repository/user_states"
	gamesearchhandlers "guru-game/internal/gamesearch/handlers"
	gamestatehandlers "guru-game/internal/gamestate/handlers"
	"guru-game/internal/recommendation"
	useractivityhandlers "guru-game/internal/useractivity/handlers"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func SetupRoutes(app *fiber.App, userStateRepo user_states.UserStateRepository, boardGameRepo boardgame.BoardGameRepository, gameRuleService *service_board.GameRuleService, gameSearchHandlers *gamesearchhandlers.GameSearchHandlers) {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Warning: .env file not found")
	}

	// Get Python service URL from environment variable, default to localhost:50051
	pythonServiceURL := os.Getenv("PYTHON_SERVICE_URL")
	if pythonServiceURL == "" {
		pythonServiceURL = "http://localhost:50051"
	}

	// --- Setup REST client for recommendation ---
	log.Println("🔌 Connecting to Python ML REST service...")
	log.Printf("🌐 Python service URL: %s", pythonServiceURL)
	restClient := recommendation.NewRESTRecommendationClient(pythonServiceURL)
	log.Println("✅ REST client initialized")

	bgService := service_board.GetBoardgameService()
	// Pass userStateRepo to the recommendation handler constructor
	recommendHandler := recommendation.NewHandler(restClient, bgService, userStateRepo)
	log.Println("✅ Recommendation handler initialized")

	// Initialize Boardgame Handlers with BoardgameRepository
	boardGameHandlers := handlers_board.NewBoardGameHandlers(boardGameRepo)

	// Auth routes
	api := app.Group("/auth")
	api.Post("/register", handlers_Auth.RegisterHandler)
	api.Post("/login", handlers_Auth.LoginHandler)
	app.Post("/auth/verify-register-otp", handlers_Auth.VerifyRegisterOTPHandler)
	app.Post("/auth/verify-login-otp", handlers_Auth.VerifyLoginOTPHandler)
	api.Post("/resend-otp", handlers_Auth.ResendOTPHandler)

	api.Get("/status", jwt.JWTMiddleware, handlers_Auth.StatusHandler)
	api.Get("/users", handlers_Auth.GetAllUsersHandler)
	api.Get("/profile", jwt.JWTMiddleware, handlers_Auth.GetProfileHandler)
	api.Put("/user/update", jwt.JWTMiddleware, handlers_Auth.UpdateUserHandler)
	api.Delete("/user/delete", jwt.JWTMiddleware, handlers_Auth.DeleteUserHandler)

	// Boardgame routes
	bg := app.Group("/boardgames")
	// Apply JWT middleware to potentially get user ID, but handler logic should handle unauthenticated users
	bg.Get("/", jwt.JWTMiddleware, boardGameHandlers.HandleGetAllBoardGames)
	bg.Get("/:id", boardGameHandlers.GetBoardGameByIDHandler)
	bg.Get("/es/:id", boardGameHandlers.GetBoardGameByIDFromESHandler)

	// User Activity routes
	userActivity := app.Group("/user/activities")
	// Create a new instance of UserActivityHandler with the restClient
	userActivityHandler := useractivityhandlers.NewUserActivityHandler(restClient)
	userActivity.Post("/", userActivityHandler.HandleUserActivity)

	// Recommendation routes
	reco := app.Group("/recommendations")

	// ส่งข้อมูล boardgames ทั้งหมดไปยัง Python ML service
	reco.Post("/send-all", recommendHandler.HandleSendAllBoardgames)
	reco.Get("/send-all", recommendHandler.HandleSendAllBoardgames)

	// ขอ recommendations สำหรับ user
	reco.Get("/", recommendHandler.HandleGetRecommendations)
	reco.Get("/user/:user_id", func(c *fiber.Ctx) error {
		c.Queries()["user_id"] = c.Params("user_id")
		return recommendHandler.HandleGetRecommendations(c)
	})

	// ดึง boardgames ทั้งหมดจาก Elasticsearch ผ่าน service
	reco.Get("/all-boardgames", recommendHandler.HandleGetAllBoardgamesFromES)

	// ขอ popular boardgames
	reco.Get("/popular", recommendHandler.HandleGetPopularBoardgames)

	// User actions
	reco.Post("/actions", recommendHandler.HandleAddUserAction)
	reco.Get("/actions/user/:user_id", recommendHandler.HandleGetUserActions)
	reco.Get("/actions/boardgame/:boardgame_id", recommendHandler.HandleGetBoardgameActions)

	// Get user's favorite boardgames directly from DB
	reco.Get("/favorites/:user_id", recommendHandler.HandleGetFavoritedBoardgames)
	// Game State Update routes
	gameState := app.Group("/api/game/updateState")
	// Create an instance of GameStateHandlers with the UserStateRepository
	gameStateHandlersInstance := gamestatehandlers.NewGameStateHandlers(userStateRepo)
	gameState.Post("/", gameStateHandlersInstance.HandleGameStateUpdate)
	gameState.Put("/", gameStateHandlersInstance.HandleGameStateUpdate)
	gameState.Patch("/", gameStateHandlersInstance.HandleGameStateUpdate)

	// Game Search routes
	gameSearch := app.Group("/api/search")
	gameSearch.Get("/", gameSearchHandlers.HandleGameSearch)
}

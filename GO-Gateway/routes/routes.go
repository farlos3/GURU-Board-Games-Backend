package routes

import (
	"log"
	"github.com/gofiber/fiber/v2"
	"guru-game/internal/auth/handlers_Auth"
	"guru-game/internal/auth/jwt"
	"guru-game/internal/boardgame/handlers_board"
	"guru-game/internal/boardgame/service_board"
	"guru-game/internal/recommendation"
)

func SetupRoutes(app *fiber.App) {
	// --- Setup gRPC client for recommendation ---
	log.Println("🔌 Connecting to Python ML gRPC service...")
	grpcClient, err := recommendation.NewGRPCRecommendationClient("localhost:8001")
	if err != nil {
		log.Printf("⚠️  Warning: Failed to connect to Python ML gRPC service: %v", err)
		log.Println("📝 Recommendation features will be disabled")
		// ไม่ panic แต่จะทำให้ recommendation routes ไม่ทำงาน
	}

	bgService := service_board.GetBoardgameService()
	
	var recommendHandler *recommendation.RecommendationHandler
	if grpcClient != nil {
		recommendHandler = recommendation.NewRecommendationHandler(grpcClient, bgService)
		log.Println("✅ Recommendation handler initialized")
	}

	// Auth routes
	api := app.Group("/auth")
	api.Post("/register", handlers_Auth.RegisterHandler)
	api.Post("/login", handlers_Auth.LoginHandler)
	app.Post("/auth/verify-register-otp", handlers_Auth.VerifyRegisterOTPHandler)
	app.Post("/auth/verify-login-otp", handlers_Auth.VerifyLoginOTPHandler)
	api.Post("/resend-otp", handlers_Auth.ResendOTPHandler)

	api.Get("/status", jwt.JWTMiddleware, handlers_Auth.StatusHandler)
	api.Get("/users", handlers_Auth.GetAllUsersHandler)
	api.Put("/user/update", jwt.JWTMiddleware, handlers_Auth.UpdateUserHandler)
	api.Delete("/user/delete", jwt.JWTMiddleware, handlers_Auth.DeleteUserHandler)

	// Boardgame routes
	bg := app.Group("/boardgames")
	bg.Get("/", handlers_board.GetAllBoardGamesHandler)
	bg.Get("/:id", handlers_board.GetBoardGameByIDHandler)
	bg.Delete("/:id", handlers_board.DeleteBoardGameHandler)

	// Recommendation routes
	reco := app.Group("/recommendations")
	
	if recommendHandler != nil {
		// ส่งข้อมูล boardgames ทั้งหมดไปยัง Python ML service
		reco.Post("/send-all", recommendHandler.HandleSendAllBoardgames)
		reco.Get("/send-all", recommendHandler.HandleSendAllBoardgames) // รองรับ GET ด้วย
		
		// ขอ recommendations สำหรับ user
		reco.Get("/", recommendHandler.HandleGetRecommendations)
		reco.Get("/user/:user_id", func(c *fiber.Ctx) error {
			// ตั้งค่า user_id จาก path parameter
			c.Queries()["user_id"] = c.Params("user_id")
			return recommendHandler.HandleGetRecommendations(c)
		})
	} else {
		// หาก gRPC client ไม่พร้อมใช้งาน
		reco.All("/*", func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error": "Recommendation service is currently unavailable",
				"message": "Python ML service is not connected",
			})
		})
	}
}
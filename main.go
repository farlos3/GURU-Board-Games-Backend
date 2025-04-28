package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"guru-game/internal/auth/service_auth"
	"guru-game/internal/boardgame/service_board"

	"guru-game/internal/db/connection"
	"guru-game/internal/db/repository/boardgame"
	"guru-game/internal/db/repository/user"
	"guru-game/routes"
)

func main() {
	// สร้างแอปพลิเคชัน Fiber
	app := fiber.New()

	// เชื่อมต่อกับฐานข้อมูล
	connection.ConnectDB()

	// ✅ Inject PostgresUserRepository เข้าไปใน auth.Init()
	service_auth.Init(&user.PostgresUserRepository{})

	// ✅ Inject PostgresBoardgameRepository เข้าไปใน service_board.Init()
	service_board.Init(&boardgame.PostgresBoardgameRepository{})

	// ตั้งค่า Routes
	routes.SetupRoutes(app)

	// เริ่มเซิร์ฟเวอร์
	log.Println("🚀 Server is running on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}

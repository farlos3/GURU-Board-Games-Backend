package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"guru-game/routes"
	"guru-game/internal/db/connection"
	"guru-game/internal/db/repository"
	"guru-game/internal/auth"            
)

func main() {
	// สร้างแอปพลิเคชัน Fiber
	app := fiber.New()

	// เชื่อมต่อกับฐานข้อมูล
	connection.ConnectDB()

	// ✅ Inject PostgresUserRepository เข้าไปใน auth.Init()
	auth.Init(&db.PostgresUserRepository{})

	// ตั้งค่า Routes
	routes.SetupRoutes(app)

	// เริ่มเซิร์ฟเวอร์
	log.Println("🚀 Server is running on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}
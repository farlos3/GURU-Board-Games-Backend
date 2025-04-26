package connection

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool

// ConnectDB ทำหน้าที่เชื่อมต่อกับฐานข้อมูล Supabase
func ConnectDB() {
	// โหลด .env ก่อนใช้งาน
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// ดึงค่า DATABASE_URL จาก environment variable
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	var err error
	DB, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Failed to create connection pool: %v", err)
	}

	err = DB.Ping(context.Background())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("✅ Connected to Supabase database successfully!")
}
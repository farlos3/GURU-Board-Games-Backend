package service_board

import (
	"encoding/json"
	"errors"
	"fmt"
	"guru-game/models"
	"log"
	"net/http"
	"os"
)

// GetBoardGameByIDFromES ดึงข้อมูลบอร์ดเกมตาม ID จาก Python service
func GetBoardGameByIDFromES(id int) (*models.BoardGame, error) {
	// รับ Python service URL จาก environment variable
	pythonServiceURL := os.Getenv("PYTHON_SERVICE_URL")
	if pythonServiceURL == "" {
		pythonServiceURL = "http://localhost:50051" // default URL
	}

	// สร้าง URL สำหรับเรียก API
	url := fmt.Sprintf("%s/api/boardgames/%d", pythonServiceURL, id)

	// ส่ง GET request ไปยัง Python service
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to send request to Python service: %v\n", err)
		return nil, errors.New("failed to connect to Python service: " + err.Error())
	}
	defer resp.Body.Close()

	// ตรวจสอบ status code
	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("boardgame not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Python service returned status code: %d", resp.StatusCode)
	}

	// แปลง response เป็น BoardGame model
	var boardgame models.BoardGame
	if err := json.NewDecoder(resp.Body).Decode(&boardgame); err != nil {
		log.Printf("Failed to decode response from Python service: %v\n", err)
		return nil, errors.New("failed to decode response: " + err.Error())
	}

	return &boardgame, nil
}

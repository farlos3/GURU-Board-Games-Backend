package db

import (
	"fmt"
	"guru-game/models"
)

// Data storage mock
var Boardgames []models.BoardGame
var BoardgameIDCounter = 1

// ConnectMock ฟังก์ชันนี้จะทำการกำหนดค่าเกมเริ่มต้นในระบบ
func ConnectMockGame() {
	Boardgames = []models.BoardGame{
		{
			ID:          1,
			Name:        "Catan",
			Description: "Trade and build in the island of Catan.",
			Players:     3,
		},
		{
			ID:          2,
			Name:        "Ticket to Ride",
			Description: "Build train routes across the USA.",
			Players:     2,
		},
	}

	// Set BoardgameIDCounter to the highest ID value from the Boardgames slice
	BoardgameIDCounter = getMaxBoardgameID()
}

// Function to find the highest Boardgame ID
func getMaxBoardgameID() int {
	maxID := 0
	for _, game := range Boardgames {
		if game.ID > maxID {
			maxID = game.ID
		}
	}
	return maxID + 1 // Increment to be ready for the next ID
}

// BoardgameRepository interface สำหรับการจัดการข้อมูลเกม
type BoardgameRepository interface {
	GetByName(name string) (*models.BoardGame, error)
	Create(boardgame *models.BoardGame) (*models.BoardGame, error)
	GetAll() []models.BoardGame
}

// MockBoardgameRepository ทำหน้าที่เป็น mock repository สำหรับ Boardgame
type MockBoardgameRepository struct{}

func (r *MockBoardgameRepository) GetByName(name string) (*models.BoardGame, error) {
	// ค้นหาเกมจาก mock data
	for _, game := range Boardgames {
		if game.Name == name {
			return &game, nil
		}
	}
	return nil, fmt.Errorf("boardgame not found")
}

func (r *MockBoardgameRepository) Create(boardgame *models.BoardGame) (*models.BoardGame, error) {
	// เพิ่มเกมใหม่เข้า mock data
	boardgame.ID = BoardgameIDCounter
	BoardgameIDCounter++
	Boardgames = append(Boardgames, *boardgame)
	return boardgame, nil
}

func (r *MockBoardgameRepository) GetAll() []models.BoardGame {
	// คืนค่าทุกเกม
	return Boardgames
}
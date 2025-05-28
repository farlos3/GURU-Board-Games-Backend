package service_board

import (
	"log"
	"errors"
	"guru-game/models"
)

// ฟังก์ชันดึงข้อมูลบอร์ดเกมทั้งหมด
func GetAllBoardGames() ([]models.BoardGame, error) {
	if boardGameRepo == nil {
		log.Println("Boardgame repository is not initialized.")
		return nil, errors.New("boardgame repository is not initialized")
	}

	log.Println("Fetching boardgames from database...")
	boardgames, err := boardGameRepo.GetAll()
	if err != nil {
		log.Printf("Failed to get boardgames: %v\n", err)
		return nil, errors.New("failed to get boardgames: " + err.Error())
	}

	if len(boardgames) == 0 {
		log.Println("No boardgames found.")
		return []models.BoardGame{}, nil
	}

	log.Printf("Successfully fetched %d boardgames.\n", len(boardgames))
	return boardgames, nil
}
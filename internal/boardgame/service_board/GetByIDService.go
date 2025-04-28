package service_board

import (
	"log"
	"errors"
	"guru-game/models"
)

// ฟังก์ชันดึงข้อมูลบอร์ดเกมตาม ID
func GetBoardGameByID(id int) (*models.BoardGame, error) {
	if boardGameRepo == nil {
		log.Println("Boardgame repository is not initialized.")
		return nil, errors.New("boardgame repository is not initialized")
	}

	boardgame, err := boardGameRepo.GetByID(id)
	if err != nil {
		log.Printf("Failed to get boardgame with ID %d: %v\n", id, err)
		return nil, errors.New("failed to get boardgame: " + err.Error())
	}

	return boardgame, nil
}
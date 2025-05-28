package service_board

import (
	"log"
	"errors"
	"guru-game/models"
)

// ฟังก์ชันสร้างบอร์ดเกมใหม่
func CreateBoardGame(boardgame *models.BoardGame) (*models.BoardGame, error) {
	if boardGameRepo == nil {
		log.Println("Boardgame repository is not initialized.")
		return nil, errors.New("boardgame repository is not initialized")
	}

	createdBoardgame, err := boardGameRepo.Create(boardgame)
	if err != nil {
		log.Printf("Failed to create boardgame: %v\n", err)
		return nil, errors.New("failed to create boardgame: " + err.Error())
	}

	return createdBoardgame, nil
}
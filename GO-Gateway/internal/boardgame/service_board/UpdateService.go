package service_board

import (
	"log"
	"errors"
	"guru-game/models"
)

// ฟังก์ชันอัปเดตบอร์ดเกม
func UpdateBoardGame(boardgame *models.BoardGame) (*models.BoardGame, error) {
	if boardGameRepo == nil {
		log.Println("Boardgame repository is not initialized.")
		return nil, errors.New("boardgame repository is not initialized")
	}

	updatedBoardgame, err := boardGameRepo.Update(boardgame)
	if err != nil {
		log.Printf("Failed to update boardgame: %v\n", err)
		return nil, errors.New("failed to update boardgame: " + err.Error())
	}

	return updatedBoardgame, nil
}
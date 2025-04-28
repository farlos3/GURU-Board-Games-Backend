package service_board

import (
	"log"
	"errors"
)

// ฟังก์ชันลบบอร์ดเกม
func DeleteBoardGame(id int) error {
	if boardGameRepo == nil {
		log.Println("Boardgame repository is not initialized.")
		return errors.New("boardgame repository is not initialized")
	}

	err := boardGameRepo.Delete(id)
	if err != nil {
		log.Printf("Failed to delete boardgame with ID %d: %v\n", id, err)
		return errors.New("failed to delete boardgame: " + err.Error())
	}

	return nil
}
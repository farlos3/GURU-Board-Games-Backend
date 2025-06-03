package service_board

import (
	"errors"
	"guru-game/internal/db/repository/boardgame"
	"guru-game/models"
	"log"
)

// ฟังก์ชันดึงข้อมูลบอร์ดเกมทั้งหมด
// Modified to accept userID and a BoardGameRepository
func GetAllBoardGames(userID int, repo boardgame.BoardGameRepository) ([]models.BoardGame, error) {
	if repo == nil {
		log.Println("Boardgame repository is not initialized.")
		return nil, errors.New("boardgame repository is not initialized")
	}

	log.Println("Fetching boardgames from database...")
	boardgames, err := repo.GetAll() // Use the passed repository
	if err != nil {
		log.Printf("Failed to get boardgames: %v\n", err)
		return nil, errors.New("failed to get boardgames: " + err.Error())
	}

	if len(boardgames) == 0 {
		log.Println("No boardgames found.")
		return []models.BoardGame{}, nil
	}

	// If user is logged in (userID > 0), fetch and include user-specific state
	if userID > 0 {
		for i := range boardgames {
			userState, err := repo.GetUserBoardgameState(userID, boardgames[i].ID)
			if err != nil {
				log.Printf("Warning: Failed to get user state for game %d: %v\n", boardgames[i].ID, err)
				// Continue without user state for this game
			} else if userState != nil {
				// Populate user-specific fields in the boardgame struct
				boardgames[i].LikedByCurrentUser = userState.Liked
				boardgames[i].FavoritedByCurrentUser = userState.Favorited
				boardgames[i].CurrentUserRating = userState.Rating
			}
		}
	}

	log.Printf("Successfully fetched %d boardgames.\n", len(boardgames))
	return boardgames, nil
}

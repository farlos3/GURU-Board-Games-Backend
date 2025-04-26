package boardgame

import (
	"errors"
	"guru-game/models"
	"guru-game/internal/db/repository"
)

var repo db.BoardgameRepository

// Init สำหรับ Inject Repository
func Init(r db.BoardgameRepository) {
	repo = r
}

func CreateNewBoardGame(newGame *models.BoardGame) (*models.BoardGame, error) {
	// ตรวจสอบว่ามีชื่อซ้ำไหม
	if game, err := repo.GetByName(newGame.Name); err == nil && game != nil {
		return nil, errors.New("boardgame already exists")
	}

	// ถ้าไม่ซ้ำ สร้างใหม่
	return repo.Create(newGame)
}

func FindBoardGameByName(name string) (*models.BoardGame, error) {
	return repo.GetByName(name)
}

func GetAllBoardgames() ([]models.BoardGame, error) {
	return repo.GetAll()
}
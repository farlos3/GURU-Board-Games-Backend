package service_board

import (
	"guru-game/internal/db/repository/boardgame"
	"guru-game/models"
)

var boardGameRepo boardgame.BoardGameRepository

// Init ฟังก์ชันใช้ในการเซ็ต repository ของบอร์ดเกม
func Init(r boardgame.BoardGameRepository) {
	boardGameRepo = r
}

type BoardgameService struct {
    repo boardgame.BoardGameRepository
}

// เพิ่ม constructor function นี้
func NewBoardgameService(repo boardgame.BoardGameRepository) *BoardgameService {
    return &BoardgameService{
        repo: repo,
    }
}

// หรือถ้าต้องการใช้ global repo
func GetBoardgameService() *BoardgameService {
    return &BoardgameService{
        repo: boardGameRepo,
    }
}

func (s *BoardgameService) GetAllBoardgames() ([]models.BoardGame, error) {
    return s.repo.GetAll()
}
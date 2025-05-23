package boardgame

import (
	"guru-game/models"
)

// BoardGameRepository interface สำหรับ CRUD บอร์ดเกม
type BoardGameRepository interface {
	GetByID(id int) (*models.BoardGame, error)
	GetAll() ([]models.BoardGame, error)
	Create(boardgame *models.BoardGame) (*models.BoardGame, error)
	Update(boardgame *models.BoardGame) (*models.BoardGame, error)
	Delete(id int) error
}

type PostgresBoardgameRepository struct{}
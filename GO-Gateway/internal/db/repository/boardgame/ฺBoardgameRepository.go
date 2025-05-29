package boardgame

import (
	"guru-game/models"
)

// BoardGameRepository interface for CRUD
type BoardGameRepository interface {
	GetByID(id int) (*models.BoardGame, error)
	GetAll() ([]models.BoardGame, error)
	Delete(id int) error
}

type PostgresBoardgameRepository struct{}
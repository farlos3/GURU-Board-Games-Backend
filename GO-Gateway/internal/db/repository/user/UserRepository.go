package user

import (
	"guru-game/models"
)

type UserRepository interface {
	GetByUsername(username string) (*models.User, error)
	GetByCredentials(username, password string) (*models.User, error)
	Create(user *models.User) (*models.User, error)
	GetAll() ([]models.User, error)
	Update(user *models.User) (*models.User, error)
	Delete(userID int64) error
	GetByEmail(username string) (*models.User, error)
	GetByID(userID int64) (*models.User, error)
}

type PostgresUserRepository struct{}

const (
	passwordPrefix = "prefix_"
	passwordSuffix = "_suffix"
)

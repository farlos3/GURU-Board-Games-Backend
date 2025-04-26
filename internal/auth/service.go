package auth

import (
	"errors"
	"guru-game/models"
	"guru-game/internal/db"
)

var repo db.MockUserRepository

func Init(r db.MockUserRepository) {
	repo = r
}

func RegisterUser(newUser *models.User) (*models.User, error) {
	// ตรวจสอบว่า username ซ้ำในฐานข้อมูลหรือไม่
	if user, err := repo.GetByUsername(newUser.Username); err == nil && user != nil {
		return nil, errors.New("username already exists")
	}
	
	// หากไม่พบ username ซ้ำในฐานข้อมูล, สร้างผู้ใช้ใหม่
	return repo.Create(newUser)
}

func LoginUser(username, password string) (*models.User, error) {
	return repo.GetByCredentials(username, password)
}

func GetAllUsers() []models.User {
	return repo.GetAll()
}

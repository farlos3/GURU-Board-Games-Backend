package service

import (
	"guru-game/internal/db/repository/user"
)

var repo user.UserRepository

// Init สำหรับ Inject Repository
func Init(r user.UserRepository) {
	repo = r
}
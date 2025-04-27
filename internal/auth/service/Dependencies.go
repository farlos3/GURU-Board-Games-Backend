package service

import (
	"guru-game/internal/db/repository"
)

var repo db.UserRepository

// Init สำหรับ Inject Repository
func Init(r db.UserRepository) {
	repo = r
}
package service_board

import (
	"guru-game/internal/db/repository/boardgame"
)

var boardGameRepo boardgame.BoardGameRepository

// Init ฟังก์ชันใช้ในการเซ็ต repository ของบอร์ดเกม
func Init(r boardgame.BoardGameRepository) {
	boardGameRepo = r
}
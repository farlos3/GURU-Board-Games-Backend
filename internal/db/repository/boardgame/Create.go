package boardgame

import (
	"context"
	"fmt"
	"guru-game/internal/db/connection"
	"guru-game/models"
)

// ฟังก์ชันเพิ่มบอร์ดเกมใหม่
func (r *PostgresBoardgameRepository) Create(boardgame *models.BoardGame) (*models.BoardGame, error) {
	// คำสั่ง SQL สำหรับการแทรกข้อมูลบอร์ดเกมใหม่
	query := `INSERT INTO boardgames (name, description, players) VALUES ($1, $2, $3) RETURNING id`
	// การแทรกข้อมูลและรับค่า ID ที่ถูกสร้างขึ้น
	err := connection.DB.QueryRow(context.Background(), query, boardgame.Name, boardgame.Description, boardgame.Players).Scan(&boardgame.ID)
	if err != nil {
		// ถ้าเกิดข้อผิดพลาดในการแทรกข้อมูล
		return nil, fmt.Errorf("failed to create board game: %v", err)
	}
	// คืนค่าบอร์ดเกมที่สร้างเสร็จแล้ว พร้อมกับค่า error (หากไม่มีข้อผิดพลาด)
	return boardgame, nil
}
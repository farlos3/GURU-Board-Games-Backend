package boardgame

import (
	"context"
	"fmt"
	"guru-game/internal/db/connection"
)

// ฟังก์ชันลบบอร์ดเกมตาม ID
func (r *PostgresBoardgameRepository) Delete(id int) error {
	query := `DELETE FROM boardgames WHERE id = $1`
	_, err := connection.DB.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("failed to delete board game: %v", err)
	}
	return nil
}
package db

import (
	"context"
	"fmt"
	"guru-game/internal/db/connection"
	"guru-game/models"
)

type BoardgameRepository interface {
	GetByName(name string) (*models.BoardGame, error)
	Create(boardgame *models.BoardGame) (*models.BoardGame, error)
	GetAll() ([]models.BoardGame, error)
}

type PostgresBoardgameRepository struct{}

func (r *PostgresBoardgameRepository) GetByName(name string) (*models.BoardGame, error) {
	query := `SELECT id, name, description, players FROM boardgames WHERE name = $1`
	row := connection.DB.QueryRow(context.Background(), query, name)

	var game models.BoardGame
	err := row.Scan(&game.ID, &game.Name, &game.Description, &game.Players)
	if err != nil {
		return nil, fmt.Errorf("boardgame not found: %v", err)
	}

	return &game, nil
}

func (r *PostgresBoardgameRepository) Create(boardgame *models.BoardGame) (*models.BoardGame, error) {
	query := `INSERT INTO boardgames (name, description, players) VALUES ($1, $2, $3) RETURNING id`
	err := connection.DB.QueryRow(context.Background(), query, boardgame.Name, boardgame.Description, boardgame.Players).Scan(&boardgame.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert boardgame: %v", err)
	}

	return boardgame, nil
}

func (r *PostgresBoardgameRepository) GetAll() ([]models.BoardGame, error) {
	query := `SELECT id, name, description, players FROM boardgames`
	rows, err := connection.DB.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch boardgames: %v", err)
	}
	defer rows.Close()

	var boardgames []models.BoardGame
	for rows.Next() {
		var game models.BoardGame
		err := rows.Scan(&game.ID, &game.Name, &game.Description, &game.Players)
		if err != nil {
			return nil, fmt.Errorf("failed to scan boardgame: %v", err)
		}
		boardgames = append(boardgames, game)
	}

	return boardgames, nil
}
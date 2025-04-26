package db

import (
	"context"
	"fmt"
	"guru-game/internal/db/connection"
	"guru-game/models"
)

type UserRepository interface {
	GetByUsername(username string) (*models.User, error)
	GetByCredentials(username, password string) (*models.User, error)
	Create(user *models.User) (*models.User, error)
	GetAll() ([]models.User, error)
}

type PostgresUserRepository struct{}

func (r *PostgresUserRepository) GetByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, password FROM users WHERE username = $1`
	row := connection.DB.QueryRow(context.Background(), query, username)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	return &user, nil
}

func (r *PostgresUserRepository) GetByCredentials(username, password string) (*models.User, error) {
	query := `SELECT id, username, password FROM users WHERE username = $1 AND password = $2`
	row := connection.DB.QueryRow(context.Background(), query, username, password)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials: %v", err)
	}

	return &user, nil
}

func (r *PostgresUserRepository) Create(user *models.User) (*models.User, error) {
	query := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`
	err := connection.DB.QueryRow(context.Background(), query, user.Username, user.Password).Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %v", err)
	}

	return user, nil
}

func (r *PostgresUserRepository) GetAll() ([]models.User, error) {
	query := `SELECT * FROM users;`
	
	rows, err := connection.DB.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %v", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		users = append(users, user)
	}

	return users, nil
}

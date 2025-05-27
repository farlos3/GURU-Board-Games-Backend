package user

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"guru-game/internal/db/connection"
	"guru-game/models"
)


func (r *PostgresUserRepository) GetByCredentials(identifier, password string) (*models.User, error) {
	var query string
	var row = connection.DB.QueryRow(context.Background(), "", "") // เตรียมตัวก่อน

	// เช็คว่าถ้า identifier มี '@' แสดงว่าเป็น email
	if strings.Contains(identifier, "@") {
		query = `SELECT id, username, password, email, full_name, avatar_url, created_at, updated_at FROM users WHERE email = $1`
		row = connection.DB.QueryRow(context.Background(), query, identifier)
	} else {
		query = `SELECT id, username, password, email, full_name, avatar_url, created_at, updated_at FROM users WHERE username = $1`
		row = connection.DB.QueryRow(context.Background(), query, identifier)
	}

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.FullName, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	// Hash Password with prefix_ + _suffix
	finalPassword := passwordPrefix + password + passwordSuffix
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(finalPassword))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials: %v", err)
	}

	return &user, nil
}
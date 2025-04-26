package db

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"guru-game/internal/db/connection"
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
}

type PostgresUserRepository struct{}

const (
	passwordPrefix = "prefix_"
	passwordSuffix = "_suffix"
)

func (r *PostgresUserRepository) Create(user *models.User) (*models.User, error) {
	// Hash Password with prefix_ + _suffix
	finalPassword := passwordPrefix + user.Password + passwordSuffix
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(finalPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	// ใช้เวลาปัจจุบันสำหรับ created_at และ updated_at
	currentTime := time.Now()

	query := `INSERT INTO users (username, password, email, full_name, avatar_url, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err = connection.DB.QueryRow(context.Background(), query, user.Username, string(hashedPassword), user.Email, user.FullName, user.AvatarURL, currentTime, currentTime).Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %v", err)
	}

	// เปลี่ยน password ใน struct เป็น hashed ด้วย
	user.Password = string(hashedPassword)
	user.CreatedAt = currentTime
	user.UpdatedAt = currentTime

	return user, nil
}

func (r *PostgresUserRepository) GetByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, password, email, full_name, avatar_url, created_at, updated_at FROM users WHERE username = $1`
	row := connection.DB.QueryRow(context.Background(), query, username)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.FullName, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	return &user, nil
}

func (r *PostgresUserRepository) GetByCredentials(username, password string) (*models.User, error) {
	query := `SELECT id, username, password, email, full_name, avatar_url, created_at, updated_at FROM users WHERE username = $1`
	row := connection.DB.QueryRow(context.Background(), query, username)

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

func (r *PostgresUserRepository) Update(user *models.User) (*models.User, error) {
	var hashedPassword string

	// เช็กว่ามีการส่ง password ใหม่มาจริงไหม
	if user.Password != "" {
		// ถ้ามี password ใหม่ ให้ทำการ hash
		finalPassword := passwordPrefix + user.Password + passwordSuffix
		hashed, err := bcrypt.GenerateFromPassword([]byte(finalPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %v", err)
		}
		hashedPassword = string(hashed)
	} else {
		// ถ้าไม่มี password ใหม่ ให้ใช้ password เดิม (ที่ถูก hash แล้ว)
		hashedPassword = user.Password
	}

	// UPDATE users
	query := `UPDATE users SET password = $1, full_name = $2, avatar_url = $3 WHERE id = $4`
	_, err := connection.DB.Exec(context.Background(), query, hashedPassword, user.FullName, user.AvatarURL, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	// ดึงข้อมูล User ใหม่จาก DB
	var updatedUser models.User
	selectQuery := `SELECT id, username, email, full_name, avatar_url, created_at, updated_at FROM users WHERE id = $1`
	err = connection.DB.QueryRow(context.Background(), selectQuery, user.ID).Scan(
		&updatedUser.ID,
		&updatedUser.Username,
		&updatedUser.Email,
		&updatedUser.FullName,
		&updatedUser.AvatarURL,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated user: %v", err)
	}

	return &updatedUser, nil
}

func (r *PostgresUserRepository) GetByEmail(email string) (*models.User, error) {
	query := `SELECT id, username, email, password, full_name, avatar_url, created_at, updated_at FROM users WHERE email = $1`
	row := connection.DB.QueryRow(context.Background(), query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.FullName, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("user not found by email: %v", err)
	}
	return &user, nil
}

func (r *PostgresUserRepository) Delete(userID int64) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := connection.DB.Exec(context.Background(), query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}
	return nil
}

func (r *PostgresUserRepository) GetAll() ([]models.User, error) {
	query := `SELECT id, username, password, email, full_name, avatar_url, created_at, updated_at FROM users`
	rows, err := connection.DB.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %v", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.FullName, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		users = append(users, user)
	}

	return users, nil
}

package user

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"guru-game/internal/db/connection"
	"guru-game/models"
)

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
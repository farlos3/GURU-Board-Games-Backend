package service_auth

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
	"guru-game/models"
)

// DeleteUser ใช้ username, email, password เพื่อลบ user
func DeleteUser(input *models.User) error {
	// หา user จาก username
	existingUser, err := repo.GetByUsername(input.Username)
	if err != nil {
		return errors.New("user not found")
	}

	// เช็ก email
	if existingUser.Email != input.Email {
		return errors.New("invalid email")
	}

	// เช็ก password (bcrypt)
	finalPassword := "prefix_" + input.Password + "_suffix"
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(finalPassword))
	if err != nil {
		return errors.New("invalid password")
	}

	// ถ้าทุกอย่างถูกต้อง ลบ user
	err = repo.Delete(existingUser.ID)
	if err != nil {
		return errors.New("failed to delete user")
	}

	log.Printf("Deleted user '%s' successfully\n", existingUser.Username)
	return nil
}
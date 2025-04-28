package service_auth

import (
	"errors"
	"log"

	"guru-game/models"
)

// UpdateUser รับ models.User ที่อาจมี username หรือ email, ทำการหา user และอัปเดต
func UpdateUser(input *models.User) (*models.User, error) {
	if input == nil {
		return nil, errors.New("input cannot be nil")
	}

	// ต้องมี username หรือ email
	if input.Username == "" && input.Email == "" {
		log.Println("Missing username or email")
		return nil, errors.New("username or email required")
	}

	// หา user ตัวจริงจาก repo
	var user *models.User
	var err error

	if input.Username != "" {
		user, err = repo.GetByUsername(input.Username)
	} else {
		user, err = repo.GetByEmail(input.Email)
	}

	if err != nil {
		log.Println("User not found:", err)
		return nil, errors.New("user not found")
	}

	log.Printf("🔎 Found user ID: %d\n", user.ID)

	// ใช้ ID จริงของ user มาเซ็ตใน input
	input.ID = user.ID

	// เช็กว่าได้ ID ไหม
	if input.ID == 0 {
		log.Println("Missing user ID after lookup")
		return nil, errors.New("missing user ID for update")
	}

	// Call repository เพื่ออัปเดต
	updatedUser, err := repo.Update(input)
	if err != nil {
		log.Printf("Failed to update user ID %d: %v\n", input.ID, err)
		return nil, err
	}

	log.Printf("User ID %d updated successfully.\n", updatedUser.ID)
	return updatedUser, nil
}
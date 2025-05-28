package service_auth

import (
	"errors"
	"log"

	"guru-game/models"
)

// UpdateUser รับ models.User ที่อาจมี username หรือ email, ทำการหา user และอัปเดต
// เพิ่มการตรวจสอบว่า user ที่มาจาก token คือผู้ใช้นั้นหรือไม่
func UpdateUser(input *models.User, tokenUserID int64) (*models.User, error) {
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

	// เช็กว่า ID ของผู้ใช้จาก token ตรงกับ ID ของ user ที่ต้องการอัปเดต
	if user.ID != tokenUserID {
		log.Println("User ID from token does not match the user being updated")
		return nil, errors.New("you can only update your own information")
	}

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
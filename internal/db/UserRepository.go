package db

import (
	"fmt"
	"guru-game/models"
)

// Data storage mock
var Users []models.User
var UserIDCounter = 1

// ConnectMock ฟังก์ชันนี้จะทำการกำหนดค่าผู้ใช้เริ่มต้นในระบบ
func ConnectMock() {
	Users = []models.User{
		{
			ID:       1,
			Username: "admin",
			Password: "admin123",
		},
		{
			ID:       2,
			Username: "testuser",
			Password: "pass1234",
		},
	}

	// Set UserIDCounter to the highest ID value from the Users slice
	UserIDCounter = getMaxUserID()
}

// Function to find the highest UserID
func getMaxUserID() int {
	maxID := 0
	for _, user := range Users {
		if user.ID > maxID {
			maxID = user.ID
		}
	}
	return maxID + 1 // Increment to be ready for the next ID
}

// UserRepository interface สำหรับการจัดการข้อมูลผู้ใช้
type UserRepository interface {
	GetByUsername(username string) (*models.User, error)
	GetByCredentials(username, password string) (*models.User, error)
	Create(user *models.User) (*models.User, error)
	GetAll() []models.User
}

// MockUserRepository ทำหน้าที่เป็น mock repository สำหรับ User
type MockUserRepository struct{}

func (r *MockUserRepository) GetByUsername(username string) (*models.User, error) {
	// การค้นหาจาก mock data
	for _, user := range Users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (r *MockUserRepository) GetByCredentials(username, password string) (*models.User, error) {
	// ตรวจสอบ username/password
	for _, user := range Users {
		if user.Username == username && user.Password == password {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("invalid credentials")
}

func (r *MockUserRepository) Create(user *models.User) (*models.User, error) {
	// เพิ่ม user เข้าไปใน mock data
	user.ID = UserIDCounter
	UserIDCounter++
	Users = append(Users, *user)
	return user, nil
}

func (r *MockUserRepository) GetAll() []models.User {
	// คืนค่าผู้ใช้ทั้งหมด
	return Users
}
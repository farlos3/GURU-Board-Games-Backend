package otp

import (
	"crypto/rand"
	"fmt"
	"time"
	"sync"
)

// OTP เก็บข้อมูลรหัสและเวลาหมดอายุ
type OTP struct {
	Code      string
	ExpiresAt time.Time
}

// จำลอง store ด้วย map ในหน่วยความจำ
var otpStore = make(map[string]OTP)
var mu sync.Mutex // ใช้ lock เพื่อความปลอดภัยในกรณีที่ใช้หลาย goroutines

// GenerateOTP สร้างรหัส OTP แบบ 6 หลัก
func GenerateOTP() (string, error) {
	b := make([]byte, 3)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// limit ให้เป็นเลข 6 หลัก
	return fmt.Sprintf("%06d", int(b[0])<<16|int(b[1])<<8|int(b[2])%1000000), nil
}

// SaveOTP บันทึกรหัส OTP พร้อมวันหมดอายุ
func SaveOTP(username, code string) {
	mu.Lock() // ใช้ Lock เพื่อให้มั่นใจว่ามีการเข้าถึงข้อมูลพร้อมกันไม่ได้
	defer mu.Unlock()

	otpStore[username] = OTP{
		Code:      code,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
}

// VerifyOTP ตรวจสอบว่า OTP ถูกต้องและยังไม่หมดอายุ
func VerifyOTP(username, code string) bool {
	mu.Lock() // ใช้ Lock เพื่อให้มั่นใจว่ามีการเข้าถึงข้อมูลพร้อมกันไม่ได้
	defer mu.Unlock()

	otp, exists := otpStore[username]
	if !exists {
		return false
	}
	if time.Now().After(otp.ExpiresAt) {
		delete(otpStore, username)
		return false
	}
	if otp.Code != code {
		return false
	}
	delete(otpStore, username) // ใช้แล้วลบทิ้ง
	return true
}
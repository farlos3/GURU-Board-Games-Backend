package otp

import (
	"fmt"
	"net/smtp"
	"os"
	"log"
)

// SendEmail ส่ง OTP ผ่าน Gmail SMTP
func SendEmail(to string, code string) error {
	from := os.Getenv("EMAIL_FROM")       
	password := os.Getenv("EMAIL_PASSWORD") 

	// ข้อความที่ต้องการส่ง
	msg := []byte("To: " + to + "\r\n" +
		"Subject: Your OTP Code\r\n" +
		"\r\n" +
		"Your OTP code is: " + code + "\r\n")

	// ตั้งค่า SMTP Auth
	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")

	// ส่งอีเมลผ่าน Gmail SMTP Server
	err := smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	// Log ว่าส่ง OTP ไปที่อีเมลแล้ว
	log.Printf("OTP code sent to %s successfully", to)
	return nil
}
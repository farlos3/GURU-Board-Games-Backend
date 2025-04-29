package jwt

import (
	"errors"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey = os.Getenv("JWT_SECRET") // ดึงค่า secret key จาก environment

// VerifyToken ตรวจสอบว่า token ถูกต้องและส่งกลับ claims
func VerifyToken(tokenString string) (*Claims, error) {
	// ตรวจสอบว่า tokenString มีค่า
	if tokenString == "" {
		return nil, errors.New("token is empty")
	}

	// ตรวจสอบว่า secret key มีค่า
	if jwtSecretKey == "" {
		return nil, errors.New("JWT_SECRET environment variable not set")
	}

	// Parse JWT token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// ตรวจสอบว่า token method เป็น HMAC และใช้ secret key ที่ตั้งไว้
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecretKey), nil
	})

	// ตรวจสอบข้อผิดพลาด
	if err != nil {
		log.Println("Error parsing token:", err)
		return nil, err
	}

	// ตรวจสอบว่า claims เป็นประเภทที่คาดหวัง
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}
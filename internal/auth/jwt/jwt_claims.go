package jwt

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
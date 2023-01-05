package models

import "github.com/golang-jwt/jwt/v4"

type JwtClaims struct {
	jwt.RegisteredClaims
	UserId   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username"`
	Hash     string `json:"hash"`
	Email    string `json:"email"`
}

package dto

import "lib/data/models/auth"

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func NewUserFromModel(user auth.User) User {
	return User{
		Username: user.Username,
		Email:    user.Email,
	}
}

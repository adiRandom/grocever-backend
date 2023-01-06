package models

import (
	"lib/data/dto/auth"
)

type Register struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func NewRegisterFromDto(register auth.RegisterRequest) Register {
	return Register{
		Username: register.Username,
		Password: register.Password,
		Email:    register.Email,
	}
}

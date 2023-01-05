package models

import "auth/data/dto"

type Register struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func NewRegisterFromDto(register dto.RegisterRequest) Register {
	return Register{
		Username: register.Username,
		Password: register.Password,
		Email:    register.Email,
	}
}

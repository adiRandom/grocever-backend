package repository

import (
	"api_gateway/data/entity"
	"api_gateway/data/models"
	"api_gateway/services"
	"lib/data/database"
	"lib/data/database/repositories"
)

type User struct {
	repositories.Repository[entity.User]
}

var userRepository *User

func GetUserRepository() *User {
	if userRepository == nil {
		userRepository = &User{}
		db, err := database.GetDb()
		if err != nil {
			panic(err)
		}
		userRepository.Db = db
	}
	return userRepository
}

func (r *User) CreateFromAuth(authModel models.Auth) (*entity.User, error) {
	hash, err := services.HashPassword(authModel.Password)
	if err != nil {
		return nil, err
	}

	user := entity.User{
		Username: authModel.Username,
		Email:    authModel.Email,
		Hash:     hash,
	}

	err = r.Db.Create(&user).Error
	return &user, err
}

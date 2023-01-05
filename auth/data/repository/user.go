package repository

import (
	"auth/data/entity"
	"auth/data/models"
	"auth/services/crypto/password"
	"gorm.io/gorm"
	"lib/data/database"
	"lib/data/database/repositories"
	"lib/data/models/auth"
	"lib/helpers"
)

type User struct {
	repositories.RepositoryWithModel[entity.User, auth.User]
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
		userRepository.ToModel = toModel
		userRepository.ToEntity = toEntity
	}
	return userRepository
}

func toEntity(model auth.User) (entity.User, error) {
	return entity.User{
		Model: gorm.Model{
			ID: model.ID,
		},
		Username: model.Username,
		Hash:     model.Hash,
		Email:    model.Email,
	}, nil
}

func toModel(entity entity.User) (auth.User, error) {
	return auth.User{
		ID:       entity.ID,
		Username: entity.Username,
		Hash:     entity.Hash,
		Email:    entity.Email,
	}, nil
}

func (r *User) CreateFromAuth(authModel models.Register) (*auth.User, error) {
	hash, err := password.HashPassword(authModel.Password)
	if err != nil {
		return nil, err
	}

	user := entity.User{
		Username: authModel.Username,
		Email:    authModel.Email,
		Hash:     hash,
	}

	err = r.Db.Create(&user).Error
	userModel, err := r.ToModel(user)
	if err != nil {
		return nil, err
	}

	return &userModel, nil
}

func (r *User) GetByUsernameAndPwd(username string, pwd string) (*auth.User, error) {
	user := entity.User{}
	err := r.Db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	if err := password.VerifyPassword(user.Hash, pwd); err != nil {
		return nil, helpers.Error{Msg: "Invalid username or password"}
	}

	userModel, err := r.ToModel(user)
	if err != nil {
		return nil, err
	}

	return &userModel, nil
}

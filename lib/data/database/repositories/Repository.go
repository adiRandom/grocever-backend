package repositories

import "gorm.io/gorm"

type repository[T any] interface {
	GetAll() ([]T, error)
	GetById(id uint) (*T, error)
	Save(entity T) error
	Delete(entity T) error
	Create(entity T) error
}

type Repository[T any] struct {
	Db *gorm.DB
}

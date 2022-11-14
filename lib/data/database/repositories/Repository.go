package repositories

import (
	"errors"
	"gorm.io/gorm"
)

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

func (r *Repository[T]) GetAll() ([]T, error) {
	var res []T
	err := r.Db.Find(&res).Error
	return res, err
}

func (r *Repository[T]) GetById(id uint) (*T, error) {
	var entity T
	err := r.Db.First(&entity, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &entity, err
}

func (r *Repository[T]) Save(entity T) error {
	return r.Db.Save(&entity).Error
}

func (r *Repository[T]) Delete(entity T) error {
	return r.Db.Delete(&entity).Error
}

func (r *Repository[T]) Create(entity T) error {
	return r.Db.Create(&entity).Error
}

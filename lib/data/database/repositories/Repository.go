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

type DbRepository[T any] struct {
	Db *gorm.DB
}

func (r *DbRepository[T]) GetAll() ([]T, error) {
	var res []T
	err := r.Db.Find(&res).Error
	return res, err
}

func (r *DbRepository[T]) GetById(id uint) (*T, error) {
	var entity T
	err := r.Db.First(&entity, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &entity, err
}

func (r *DbRepository[T]) Save(entity T) error {
	return r.Db.Save(&entity).Error
}

func (r *DbRepository[T]) Delete(entity T) error {
	return r.Db.Delete(&entity).Error
}

func (r *DbRepository[T]) Create(entity *T) error {
	return r.Db.Create(entity).Error
}

func (r *DbRepository[T]) CreateMany(entities []T) error {
	return r.Db.Create(&entities).Error
}

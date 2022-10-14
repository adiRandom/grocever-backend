package repositories

type Repository[T any] interface {
	GetAll() ([]T, error)
	GetById(id uint) (*T, error)
	Save(entity T) error
	Update(entity T) error
}

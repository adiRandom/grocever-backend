package repositories

import "productProcessing/data/database/entities"

type ProductRepository struct {
	Repository[entities.ProductEntity]
}

var pr *ProductRepository = nil

func GetProductRepository() *ProductRepository {
	if pr == nil {
		pr = &ProductRepository{}
	}
	return pr
}

func (r *ProductRepository) GetAll() ([]entities.ProductEntity, error) {
	var products []entities.ProductEntity
	err := r.db.Find(&products).Error
	return products, err
}

func (r *ProductRepository) GetById(id uint) (*entities.ProductEntity, error) {
	var product entities.ProductEntity
	err := r.db.First(&product, id).Error
	return &product, err
}

func (r *ProductRepository) Save(entity entities.ProductEntity) error {
	return r.db.Save(&entity).Error
}

func (r *ProductRepository) Delete(entity entities.ProductEntity) error {
	return r.db.Delete(&entity).Error
}

func (r *ProductRepository) GetProductByNameAndStoreId(name string, storeId int) (*entities.ProductEntity, error) {
	var product entities.ProductEntity
	err := r.db.Where("name = ? AND store_id = ?", name, storeId).First(&product).Error
	return &product, err
}

package repositories

import (
	"dealScraper/data/database"
	"dealScraper/data/database/entities"
	"gorm.io/gorm"
)

type ProductWithBestOfferRepository struct {
	db *gorm.DB
}

var productWithBestOfferRepository *ProductWithBestOfferRepository = nil

func GetProductWithBestOfferRepository() (*ProductWithBestOfferRepository, error) {
	if productWithBestOfferRepository == nil {
		db, err := database.GetDb()

		if err != nil {
			return nil, err
		}
		productWithBestOfferRepository = &ProductWithBestOfferRepository{db}
	}

	return productWithBestOfferRepository, nil
}

func (r *ProductWithBestOfferRepository) GetAll() ([]entities.ProductWithBestOfferEntity, error) {
	var products []entities.ProductWithBestOfferEntity
	err := r.db.Find(&products).Error
	return products, err
}

func (r *ProductWithBestOfferRepository) GetById(id uint) (*entities.ProductWithBestOfferEntity, error) {
	var product entities.ProductWithBestOfferEntity
	err := r.db.First(&product, id).Error
	return &product, err
}

func (r *ProductWithBestOfferRepository) Save(entity entities.ProductWithBestOfferEntity) error {
	return r.db.Create(&entity).Error
}

func (r *ProductWithBestOfferRepository) Update(entity entities.ProductWithBestOfferEntity) error {
	return r.db.Save(&entity).Error
}

func (r *ProductWithBestOfferRepository) GetByProductName(productName string) (*entities.ProductWithBestOfferEntity, error) {
	var product entities.ProductWithBestOfferEntity
	err := r.db.First(&product, entities.ProductWithBestOfferEntity{Name: productName}).Error
	return &product, err
}

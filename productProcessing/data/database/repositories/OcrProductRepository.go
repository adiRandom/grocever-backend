package repositories

import (
	"errors"
	"gorm.io/gorm"
	"lib/data/database"
	"lib/data/database/repositories"
	"productProcessing/data/database/entities"
)

type OcrProductRepository struct {
	repositories.Repository[entities.OcrProductEntity]
}

var ocrRepo *OcrProductRepository = nil

func GetOcrProductRepository() *OcrProductRepository {
	if ocrRepo == nil {
		ocrRepo = &OcrProductRepository{}
		db, err := database.GetDb()
		if err != nil {
			panic(err)
		}
		ocrRepo.Db = db
	}
	return ocrRepo
}

func (r *OcrProductRepository) GetAll() ([]entities.OcrProductEntity, error) {
	var ocrProducts []entities.OcrProductEntity
	err := r.Db.Find(&ocrProducts).Error
	return ocrProducts, err
}

func (r *OcrProductRepository) GetById(id uint) (*entities.OcrProductEntity, error) {
	var ocrProduct entities.OcrProductEntity
	err := r.Db.First(&ocrProduct, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ocrProduct, err
}

func (r *OcrProductRepository) Save(entity entities.OcrProductEntity) error {
	return r.Db.Save(&entity).Error
}

func (r *OcrProductRepository) Delete(entity entities.OcrProductEntity) error {
	return r.Db.Delete(&entity).Error
}

func (r *OcrProductRepository) Create(entity entities.OcrProductEntity) error {
	return r.Db.Create(&entity).Error
}

func (r *OcrProductRepository) AddOcrProductToProduct(
	ocrProduct entities.OcrProductEntity,
	product entities.ProductEntity,
) error {
	var existingOcrProducts []entities.OcrProductEntity
	err := r.Db.Model(&product).Association("OcrProducts").Find(&existingOcrProducts)

	err = r.Db.Model(&product).Association("OcrProducts").Append(&ocrProduct)
	if err != nil {
		return err
	}

	// Link this ocr product to the eixsting ocr products
	// Then link the existing ocr products to this ocr product
	for _, existingOcrProduct := range existingOcrProducts {
		err = r.Db.Model(&ocrProduct).Association("Related").Append(&existingOcrProduct)
		if err != nil {
			return err
		}
		err = r.Db.Model(&existingOcrProduct).Association("Related").Append(&ocrProduct)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *OcrProductRepository) GetBestPrice(ocrName string) (*float32, error) {
	var ocrProduct entities.OcrProductEntity
	err := r.Db.First(&ocrProduct, "ocr_product_name = ?", ocrName).Error
	if err != nil {
		return nil, err
	}

	var relatedOcrProducts []entities.OcrProductEntity
	err = r.Db.Model(&ocrProduct).Association("Related").Find(&relatedOcrProducts)
	if err != nil {
		return nil, err
	}

	var ocrProductNames = make([]string, len(relatedOcrProducts)+1)
	ocrProductNames[0] = ocrName

	for i := 1; i < len(relatedOcrProducts); i++ {
		ocrProductNames[i] = relatedOcrProducts[i-1].OcrProductName
	}

	var bestProduct entities.ProductEntity
	err = r.Db.Where("name IN (?)", ocrProductNames).Order("price").First(&bestProduct).Error

	if err != nil {
		return nil, err
	}

	return &bestProduct.Price, nil
}

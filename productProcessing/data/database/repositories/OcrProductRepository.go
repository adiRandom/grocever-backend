package repositories

import (
	"errors"
	"gorm.io/gorm"
	"productProcessing/data/database"
	"productProcessing/data/database/entities"
)

type OcrProductRepository struct {
	Repository[entities.OcrProductEntity]
}

var ocrRepo *OcrProductRepository = nil

func GetOcrProductRepository() *OcrProductRepository {
	if ocrRepo == nil {
		ocrRepo = &OcrProductRepository{}
		db, err := database.GetDb()
		if err != nil {
			panic(err)
		}
		ocrRepo.db = db
	}
	return ocrRepo
}

func (r *OcrProductRepository) GetAll() ([]entities.OcrProductEntity, error) {
	var ocrProducts []entities.OcrProductEntity
	err := r.db.Find(&ocrProducts).Error
	return ocrProducts, err
}

func (r *OcrProductRepository) GetById(id uint) (*entities.OcrProductEntity, error) {
	var ocrProduct entities.OcrProductEntity
	err := r.db.First(&ocrProduct, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ocrProduct, err
}

func (r *OcrProductRepository) Save(entity entities.OcrProductEntity) error {
	return r.db.Save(&entity).Error
}

func (r *OcrProductRepository) Delete(entity entities.OcrProductEntity) error {
	return r.db.Delete(&entity).Error
}

func (r *OcrProductRepository) Create(entity entities.OcrProductEntity) error {
	return r.db.Create(&entity).Error
}

func (r *OcrProductRepository) AddOcrProductToProduct(
	ocrProduct entities.OcrProductEntity,
	product entities.ProductEntity,
) error {
	var existingOcrProducts []entities.OcrProductEntity
	err := r.db.Model(&product).Association("OcrProducts").Find(&existingOcrProducts)

	err = r.db.Model(&product).Association("OcrProducts").Append(&ocrProduct)
	if err != nil {
		return err
	}

	// Link this ocr product to the eixsting ocr products
	// Then link the existing ocr products to this ocr product
	for _, existingOcrProduct := range existingOcrProducts {
		err = r.db.Model(&ocrProduct).Association("Related").Append(&existingOcrProduct)
		if err != nil {
			return err
		}
		err = r.db.Model(&existingOcrProduct).Association("Related").Append(&ocrProduct)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *OcrProductRepository) GetBestPrice(ocrName string) (*float32, error) {
	var ocrProduct entities.OcrProductEntity
	err := r.db.First(&ocrProduct, "ocr_product_name = ?", ocrName).Error
	if err != nil {
		return nil, err
	}

	var relatedOcrProducts []entities.OcrProductEntity
	err = r.db.Model(&ocrProduct).Association("Related").Find(&relatedOcrProducts)
	if err != nil {
		return nil, err
	}

	var ocrProductNames = make([]string, len(relatedOcrProducts)+1)
	ocrProductNames[0] = ocrName

	for i := 1; i < len(relatedOcrProducts); i++ {
		ocrProductNames[i] = relatedOcrProducts[i-1].OcrProductName
	}

	var bestProduct entities.ProductEntity
	err = r.db.Where("name IN (?)", ocrProductNames).Order("price").First(&bestProduct).Error

	if err != nil {
		return nil, err
	}

	return &bestProduct.Price, nil
}

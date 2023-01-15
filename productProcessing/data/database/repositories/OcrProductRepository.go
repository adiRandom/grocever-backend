package repositories

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"lib/data/database"
	"lib/data/database/repositories"
	"lib/data/models/product"
	"productProcessing/data/database/entities"
)

type OcrProductRepository struct {
	repositories.DbRepository[entities.OcrProductEntity]
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

func (r *OcrProductRepository) GetById(id string) (*entities.OcrProductEntity, error) {
	var ocrProduct entities.OcrProductEntity
	err := r.Db.First(&ocrProduct, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ocrProduct, err
}

func (r *OcrProductRepository) GetByIdWithJoins(name string) (*entities.OcrProductEntity, error) {
	var ocrProduct entities.OcrProductEntity
	err := r.Db.Joins("Related").Joins("Products").Find(&ocrProduct, name).Error

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

func (r *OcrProductRepository) Create(model product.OcrProductModel) error {
	entity := entities.OcrProductEntity{
		OcrProductName: model.OcrProductName,
		BestPrice:      model.BestPrice,
	}
	return r.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(&entity).Error
}

func (r *OcrProductRepository) LinkProductAndOcrProduct(
	ocrProductName string,
	product entities.ProductEntity,
) error {
	var ocrProduct entities.OcrProductEntity

	err := r.Db.First(&ocrProduct, ocrProductName).Error
	if err != nil {
		return err
	}

	var existingOcrProducts []entities.OcrProductEntity
	err = r.Db.Model(&product).Association("OcrProducts").Find(&existingOcrProducts)

	err = r.Db.Model(&product).Association("OcrProducts").Append(&ocrProduct)
	if err != nil {
		return err
	}

	err = r.Db.Model(&ocrProduct).Association("Products").Append(&product)
	if err != nil {
		return err
	}

	// Link this ocr product to the eixsting ocr product
	// Then link the existing ocr product to this ocr product
	for _, existingOcrProduct := range existingOcrProducts {
		if existingOcrProduct.OcrProductName != ocrProductName {
			err = r.Db.Model(&ocrProduct).Association("Related").Append(&existingOcrProduct)
			if err != nil {
				return err
			}
			err = r.Db.Model(&existingOcrProduct).Association("Related").Append(&ocrProduct)
			if err != nil {
				return err
			}
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

func (r *OcrProductRepository) Exists(ocrName string) (bool, error) {
	var ocrProduct entities.OcrProductEntity
	err := r.Db.First(&ocrProduct, "ocr_product_name = ?", ocrName).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

// Function that gets an array of ocr names and returns a bool array to represent which ocr names exist in the database
func (r *OcrProductRepository) ExistsMultiple(ocrNames []string) ([]bool, error) {
	var ocrProducts []entities.OcrProductEntity
	err := r.Db.Where("ocr_product_name IN (?)", ocrNames).Find(&ocrProducts).Error
	if err != nil {
		return nil, err
	}

	var exists = make([]bool, len(ocrNames))
	for _, ocrProduct := range ocrProducts {
		for i, ocrName := range ocrNames {
			if ocrProduct.OcrProductName == ocrName {
				exists[i] = true
			}
		}
	}

	return exists, nil
}

func (r *OcrProductRepository) UpdateBestPrice(ocrName string) []error {
	ocrProduct, err := r.GetByIdWithJoins(ocrName)
	if err != nil {
		return []error{err}
	}

	// Get best price from products
	var bestPrice float32
	for _, productEntity := range ocrProduct.Products {
		if bestPrice == 0 || productEntity.Price < bestPrice {
			bestPrice = productEntity.Price
		}
	}

	if ocrProduct.BestPrice != bestPrice {
		err = r.Db.Model(&ocrProduct).Update("best_price", bestPrice).Error
		if err != nil {
			return []error{err}
		}

		errList := make([]error, 0)

		// Update best price for related ocr products
		for _, relatedOcrProduct := range ocrProduct.Related {
			err = r.Db.Model(&relatedOcrProduct).Update("best_price", bestPrice).Error
			if err != nil {
				errList = append(errList, err)
			}
		}

		if len(errList) > 0 {
			return errList
		}
	}

	return nil
}

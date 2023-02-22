package entities

import (
	"gorm.io/gorm"
	"lib/data/models/product"
	"lib/functional"
	"time"
)

// OcrProductEntity Link between product and ocr names
type OcrProductEntity struct {
	OcrProductName string `gorm:"primaryKey"`
	BestProductID  *uint
	BestProduct    *ProductEntity `gorm:"foreignKey:BestProductID;"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt      `gorm:"index"`
	Products       []*ProductEntity    `gorm:"many2many:ocr-product_product;"`
	Related        []*OcrProductEntity `gorm:"many2many:ocr-product_related;"`
}

func (entity OcrProductEntity) ToModel(withProducts bool, withRelated bool) product.OcrProductModel {
	bestProductModel := entity.BestProduct.ToModel()
	ocrProductModel := product.OcrProductModel{
		OcrProductName: entity.OcrProductName,
		BestProduct:    &bestProductModel,
	}

	if withProducts {
		ocrProductModel.Products = functional.Map(entity.Products,
			func(productEntity *ProductEntity) *product.Model {
				model := productEntity.ToModel()
				return &model
			},
		)
	}

	if withRelated {
		ocrProductModel.Related = functional.Map(entity.Related,
			func(relatedEntity *OcrProductEntity) *product.OcrProductModel {
				model := relatedEntity.ToModel(false, false)
				return &model
			},
		)
	}

	return ocrProductModel
}

func NewOcrProductEntityFromModel(model product.OcrProductModel) OcrProductEntity {
	var bestProductEntity *ProductEntity = nil
	var bestProductID *uint = nil
	if model.BestProduct != nil {
		bestProductEntity = NewProductEntityFromModel(*model.BestProduct)
		bestProductID = &bestProductEntity.ID
	}
	return OcrProductEntity{
		OcrProductName: model.OcrProductName,
		BestProduct:    bestProductEntity,
		BestProductID:  bestProductID,
	}
}

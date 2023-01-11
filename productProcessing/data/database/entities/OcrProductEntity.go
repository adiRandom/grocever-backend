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
	BestPrice      float32
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt      `gorm:"index"`
	Products       []*ProductEntity    `gorm:"many2many:ocr-product_product;"`
	Related        []*OcrProductEntity `gorm:"many2many:ocr-product_related;"`
}

func (entity OcrProductEntity) ToModel(withProducts bool, withRelated bool) product.OcrProductModel {
	ocrProductModel := product.OcrProductModel{
		OcrProductName: entity.OcrProductName,
		BestPrice:      entity.BestPrice,
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
	return OcrProductEntity{
		OcrProductName: model.OcrProductName,
		BestPrice:      model.BestPrice,
	}
}

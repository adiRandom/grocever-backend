package entities

import (
	"gorm.io/gorm"
	"productProcessing/data/models"
)

type MissLink struct {
	gorm.Model
	ProductIdFk      uint
	Product          *ProductEntity `gorm:"foreignKey:ProductIdFk;references:ID"`
	OcrProductNameFk string
	OcrProduct       *OcrProductEntity `gorm:"foreignKey:OcrProductNameFk;references:OcrProductName"`
	UserId           uint              `gorm:"unique"`
}

func (entity *MissLink) ToModel() *models.MissLink {
	productModel := entity.Product.ToModel()
	ocrProductModel := entity.OcrProduct.ToModel(false, false)

	return &models.MissLink{
		Id:         entity.ID,
		Product:    &productModel,
		OcrProduct: &ocrProductModel,
		UserId:     entity.UserId,
	}
}

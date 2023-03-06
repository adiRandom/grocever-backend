package models

import (
	"lib/data/models/product"
	"productProcessing/data/database/entities"
)

type MissLink struct {
	Id         uint
	Product    *product.Model
	OcrProduct *product.OcrProductModel
	UserId     uint
}

func (model *MissLink) ToEntity() *entities.MissLink {
	productEntity := entities.NewProductEntityFromModel(*model.Product)
	ocrProductEntity := entities.NewOcrProductEntityFromModel(*model.OcrProduct)

	entity := entities.MissLink{
		ProductIdFk:      uint(model.Product.ID),
		OcrProductNameFk: model.OcrProduct.OcrProductName,
		UserId:           model.UserId,
		Product:          productEntity,
		OcrProduct:       &ocrProductEntity,
	}

	if model.Id != 0 {
		entity.ID = model.Id
	}

	return &entity
}

func NewMissLinkModelFromEntity(entity *entities.MissLink) *MissLink {
	productModel := entity.Product.ToModel()
	ocrProductModel := entity.OcrProduct.ToModel(false, false)

	return &MissLink{
		Id:         entity.ID,
		Product:    &productModel,
		OcrProduct: &ocrProductModel,
		UserId:     entity.UserId,
	}
}
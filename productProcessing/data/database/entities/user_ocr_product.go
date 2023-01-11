package entities

import (
	"gorm.io/gorm"
	"lib/data/models"
	"lib/data/models/product"
)

type UserOcrProduct struct {
	gorm.Model
	UserId         uint
	OcrProductName string
	OcrProduct     OcrProductEntity `gorm:"foreignKey:OcrProductName"`
	Qty            float32
	UnitPrice      float32
	Price          float32
	StoreId        uint
	UnitType       string
}

func (entity UserOcrProduct) ToModel(store models.StoreMetadata) product.UserOcrProductModel {
	return product.UserOcrProductModel{
		Id:         entity.ID,
		UserId:     entity.UserId,
		OcrProduct: entity.OcrProduct.ToModel(false, false),
		Qty:        entity.Qty,
		UnitPrice:  entity.UnitPrice,
		Price:      entity.Price,
		StoreId:    entity.StoreId,
		UnitType:   entity.UnitType,
	}
}

func NewUserOcrProductFromModel(model product.UserOcrProductModel) UserOcrProduct {
	return UserOcrProduct{
		Model:          gorm.Model{ID: model.Id},
		UserId:         model.UserId,
		OcrProductName: model.OcrProduct.OcrProductName,
		Qty:            model.Qty,
		OcrProduct:     NewOcrProductEntityFromModel(model.OcrProduct),
		UnitPrice:      model.UnitPrice,
		StoreId:        model.StoreId,
		Price:          model.Price,
		UnitType:       model.UnitType,
	}
}

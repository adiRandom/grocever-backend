package entities

import (
	"gorm.io/gorm"
	"lib/data/models"
	"lib/data/models/product"
	"lib/helpers"
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
		Id:         int(entity.ID),
		UserId:     int(entity.UserId),
		OcrProduct: entity.OcrProduct.ToModel(false, false),
		Qty:        entity.Qty,
		UnitPrice:  entity.UnitPrice,
		Price:      entity.Price,
		Store:      store,
		UnitType:   entity.UnitType,
	}
}

func NewUserOcrProductFromModel(model product.UserOcrProductModel) (*UserOcrProduct, error) {
	if model.Id == -1 {
		return nil, helpers.Error{Msg: "Id is required to create user ocr product entity"}
	}

	return &UserOcrProduct{
		Model:          gorm.Model{ID: uint(model.Id)},
		UserId:         uint(model.UserId),
		OcrProductName: model.OcrProduct.OcrProductName,
		Qty:            model.Qty,
		OcrProduct:     NewOcrProductEntityFromModel(model.OcrProduct),
		UnitPrice:      model.UnitPrice,
		StoreId:        uint(model.Store.StoreId),
		Price:          model.Price,
		UnitType:       model.UnitType,
	}, nil
}

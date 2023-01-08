package entities

import (
	"gorm.io/gorm"
	"lib/data/models"
	"lib/data/models/product"
)

type UserOcrProduct struct {
	gorm.Model
	UserId         string
	OcrProductName string
	OcrProduct     OcrProductEntity `gorm:"foreignKey:OcrProductName"`
	ProductID      uint
	Product        ProductEntity `gorm:"foreignKey:ProductID"`
	Qty            float32
}

func (entity UserOcrProduct) ToModel(store models.StoreMetadata) product.UserOcrProductModel {
	return product.UserOcrProductModel{
		Id:         entity.ID,
		UserId:     entity.UserId,
		OcrProduct: entity.OcrProduct.ToModel(false, false),
		Product:    entity.Product.ToModel(),
		Qty:        entity.Qty,
		Store:      store,
		Price:      entity.Qty * entity.Product.Price,
	}
}

func NewUserOcrProductFromModel(model product.UserOcrProductModel) UserOcrProduct {
	return UserOcrProduct{
		Model:          gorm.Model{ID: model.Id},
		UserId:         model.UserId,
		OcrProductName: model.OcrProduct.OcrProductName,
		ProductID:      model.Product.ID,
		Qty:            model.Qty,
		OcrProduct:     NewOcrProductEntityFromModel(model.OcrProduct),
		Product:        NewProductEntityFromModel(model.Product),
	}
}

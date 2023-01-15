package product

import (
	"lib/data/dto/product"
	"lib/data/models"
)

type UserOcrProductModel struct {
	Id         uint
	Qty        float32
	Price      float32
	UserId     uint
	OcrProduct OcrProductModel
	UnitPrice  float32
	Store      models.StoreMetadata
	UnitType   string
}

func NewUserOcrProductModel(id uint, qty float32, price float32, userId uint, ocrProduct OcrProductModel, unitPrice float32, store models.StoreMetadata, unitType string) *UserOcrProductModel {
	return &UserOcrProductModel{Id: id, Qty: qty, Price: price, UserId: userId, OcrProduct: ocrProduct, UnitPrice: unitPrice, Store: store, UnitType: unitType}
}
func (m *UserOcrProductModel) ToDto() product.UserOcrProductDto {
	return product.UserOcrProductDto{
		Id:        m.Id,
		OcrName:   m.OcrProduct.OcrProductName,
		Qty:       m.Qty,
		UnitPrice: m.UnitPrice,
		UnitName:  m.UnitType,
		Price:     m.Price,
		Store:     m.Store,
		UserId:    m.UserId,
	}
}

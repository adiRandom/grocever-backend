package product

import (
	"lib/data/dto/product"
	"lib/data/models"
)

type PurchaseInstalmentModel struct {
	Id         int
	Qty        float32
	Price      float32
	UserId     int
	OcrProduct OcrProductModel
	UnitPrice  float32
	Store      models.StoreMetadata
	UnitType   string
}

func NewPurchaseInstalmentModel(id int, qty float32, price float32, userId int, ocrProduct OcrProductModel, unitPrice float32, store models.StoreMetadata, unitType string) *PurchaseInstalmentModel {
	return &PurchaseInstalmentModel{Id: id, Qty: qty, Price: price, UserId: userId, OcrProduct: ocrProduct, UnitPrice: unitPrice, Store: store, UnitType: unitType}
}
func (m *PurchaseInstalmentModel) ToDto() product.PurchaseInstalmentDto {
	return product.PurchaseInstalmentDto{
		Id:        m.Id,
		OcrName:   m.OcrProduct.OcrProductName,
		Qty:       m.Qty,
		UnitPrice: m.UnitPrice,
		UnitName:  m.UnitType,
		Price:     m.Price,
		Store:     m.Store.ToDto(),
	}
}

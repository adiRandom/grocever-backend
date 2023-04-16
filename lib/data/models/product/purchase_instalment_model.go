package product

import (
	"lib/data/dto/product"
	"lib/data/models"
	"time"
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
	Date       *time.Time
}

func NewPurchaseInstalmentModel(id int, qty float32, price float32, userId int, ocrProduct OcrProductModel, unitPrice float32, store models.StoreMetadata, unitType string, date *time.Time) *PurchaseInstalmentModel {
	return &PurchaseInstalmentModel{Id: id, Qty: qty, Price: price, UserId: userId, OcrProduct: ocrProduct, UnitPrice: unitPrice, Store: store, UnitType: unitType, Date: date}
}
func (m *PurchaseInstalmentModel) ToDto() product.PurchaseInstalmentDto {
	var date int64
	if m.Date != nil {
		date = m.Date.Unix()
	}
	return product.PurchaseInstalmentDto{
		Id:        m.Id,
		OcrName:   m.OcrProduct.OcrProductName,
		Qty:       m.Qty,
		UnitPrice: m.UnitPrice,
		UnitName:  m.UnitType,
		Price:     m.Price,
		Store:     m.Store.ToDto(),
		Date:      date,
	}
}

func (m *PurchaseInstalmentModel) ToCreateDto() product.CreatePurchaseInstalmentDto {
	return product.CreatePurchaseInstalmentDto{
		OcrName:   m.OcrProduct.OcrProductName,
		Qty:       m.Qty,
		UnitPrice: m.UnitPrice,
		UnitName:  m.UnitType,
		Store:     m.Store.ToDto(),
		UserId:    uint(m.UserId),
		Date:      m.Date,
	}
}

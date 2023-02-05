package entities

import (
	"gorm.io/gorm"
	"lib/data/models"
	"lib/data/models/product"
)

type PurchaseInstalment struct {
	gorm.Model
	UserId           uint
	OcrProductNameFk string           `gorm:"size:255"`
	OcrProduct       OcrProductEntity `gorm:"foreignKey:OcrProductNameFk;references:OcrProductName"`
	Qty              float32
	UnitPrice        float32
	Price            float32
	StoreId          uint
	UnitType         string
}

func NewPurchaseInstalment(
	userId uint,
	ocrProductNameFk string,
	ocrProduct OcrProductEntity,
	qty float32,
	unitPrice float32,
	price float32,
	storeId uint,
	unitType string,
) *PurchaseInstalment {
	return &PurchaseInstalment{
		UserId:           userId,
		OcrProductNameFk: ocrProductNameFk,
		OcrProduct:       ocrProduct,
		Qty:              qty,
		UnitPrice:        unitPrice,
		Price:            price,
		StoreId:          storeId,
		UnitType:         unitType,
	}
}

func (entity PurchaseInstalment) ToModel(store models.StoreMetadata) product.PurchaseInstalmentModel {
	return product.PurchaseInstalmentModel{
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

func NewPurchaseInstalmentFromModel(model product.PurchaseInstalmentModel) *PurchaseInstalment {
	entity := &PurchaseInstalment{
		UserId:           uint(model.UserId),
		OcrProductNameFk: model.OcrProduct.OcrProductName,
		Qty:              model.Qty,
		OcrProduct:       NewOcrProductEntityFromModel(model.OcrProduct),
		UnitPrice:        model.UnitPrice,
		StoreId:          uint(model.Store.StoreId),
		Price:            model.Price,
		UnitType:         model.UnitType,
	}

	if model.Id != -1 {
		entity.Model = gorm.Model{ID: uint(model.Id)}
	}

	return entity
}

package models

import (
	"gorm.io/gorm"
	"lib/data/dto/product"
	"productProcessing/data/database/entities"
)

type UpdatePurchaseInstalmentModel struct {
	Id        uint    `json:"id"`
	OcrName   string  `json:"ocrName"`
	Qty       float32 `json:"qty"`
	UnitPrice float32 `json:"unitPrice"`
	UnitName  string  `json:"unitName"`
	StoreId   uint    `json:"storeId"`
	UserId    uint    `json:"userId"`
}

func (m *UpdatePurchaseInstalmentModel) ToEntity() *entities.PurchaseInstalment {
	return &entities.PurchaseInstalment{
		Model: gorm.Model{
			ID: m.Id,
		},
		Qty:              m.Qty,
		Price:            m.Qty * m.UnitPrice,
		UnitPrice:        m.UnitPrice,
		UnitType:         m.UnitName,
		UserId:           m.UserId,
		OcrProductNameFk: m.OcrName,
		StoreId:          m.StoreId,
	}
}

func NewUpdatePurchaseInstalmentModelFromDto(
	dto *product.UpdatePurchaseInstalmentDto,
	userId uint,
	purchaseInstalmentId uint,
) *UpdatePurchaseInstalmentModel {
	return &UpdatePurchaseInstalmentModel{
		Id:        purchaseInstalmentId,
		OcrName:   dto.OcrName,
		Qty:       dto.Qty,
		UnitPrice: dto.UnitPrice,
		UnitName:  dto.UnitName,
		StoreId:   dto.StoreId,
		UserId:    userId,
	}
}

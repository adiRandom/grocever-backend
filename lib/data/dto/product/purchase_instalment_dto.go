package product

import (
	"lib/data/dto/store"
	"lib/data/models/product"
)

type PurchaseInstalmentDto struct {
	Id        int               `json:"id"`
	OcrName   string            `json:"ocrName"`
	Qty       float32           `json:"qty"`
	UnitPrice float32           `json:"unitPrice"`
	UnitName  string            `json:"unitName"`
	Price     float32           `json:"price"`
	Store     store.MetadataDto `json:"store"`
}

type PurchaseInstalmentWithUserDto struct {
	PurchaseInstalmentDto
	UserId int `json:"userId"`
}

type CretePurchaseInstalmentDto struct {
	OcrName   string            `json:"ocrName"`
	Qty       float32           `json:"qty"`
	UnitPrice float32           `json:"unitPrice"`
	UnitName  string            `json:"unitName"`
	Store     store.MetadataDto `json:"store"`
	UserId    uint              `json:"userId"`
}

func NewCreatePurchaseInstalmentDtoFromModel(
	model product.PurchaseInstalmentModel,
) CretePurchaseInstalmentDto {
	return CretePurchaseInstalmentDto{
		OcrName:   model.OcrProduct.OcrProductName,
		Qty:       model.Qty,
		UnitPrice: model.UnitPrice,
		UnitName:  model.UnitType,
		Store:     model.Store.ToDto(),
		UserId:    uint(model.UserId),
	}
}

type CreatePurchaseInstalmentListDto struct {
	Instalments []CretePurchaseInstalmentDto `json:"instalments"`
}

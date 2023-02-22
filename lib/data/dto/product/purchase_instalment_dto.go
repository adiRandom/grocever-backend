package product

import (
	"lib/data/dto/store"
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

type CreatePurchaseInstalmentDto struct {
	OcrName   string            `json:"ocrName"`
	Qty       float32           `json:"qty"`
	UnitPrice float32           `json:"unitPrice"`
	UnitName  string            `json:"unitName"`
	Store     store.MetadataDto `json:"store"`
	UserId    uint              `json:"userId"`
}

type CreatePurchaseInstalmentNoOcrDto struct {
	ProductName string  `json:"ocrName"`
	Qty         float32 `json:"qty"`
	UnitPrice   float32 `json:"unitPrice"`
	UnitName    string  `json:"unitName"`
	StoreId     uint    `json:"storeId"`
}

type CreatePurchaseInstalmentNoOcrWithUserDto struct {
	CreatePurchaseInstalmentNoOcrDto
	UserId uint `json:"userId"`
}

type CreatePurchaseInstalmentListDto struct {
	Instalments []CreatePurchaseInstalmentDto `json:"instalments"`
}

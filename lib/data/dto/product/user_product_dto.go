package product

import (
	"lib/data/dto/store"
)

type UserOcrProductDto struct {
	Id        int               `json:"id"`
	OcrName   string            `json:"ocrName"`
	Qty       float32           `json:"qty"`
	UnitPrice float32           `json:"unitPrice"`
	UnitName  string            `json:"unitName"`
	Price     float32           `json:"price"`
	Store     store.MetadataDto `json:"store"`
	UserId    int               `json:"userId"`
}

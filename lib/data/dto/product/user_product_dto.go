package product

import "lib/data/models"

type UserOcrProductDto struct {
	Id        uint                 `json:"id"`
	OcrName   string               `json:"ocrName"`
	Qty       float32              `json:"qty"`
	UnitPrice float32              `json:"unitPrice"`
	UnitName  string               `json:"unitName"`
	Price     float32              `json:"price"`
	Store     models.StoreMetadata `json:"store"`
	UserId    uint                 `json:"userId"`
}

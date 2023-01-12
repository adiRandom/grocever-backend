package models

import (
	"lib/data/dto"
	"lib/data/models"
)

type OcrProduct struct {
	Name      string
	Price     float32
	UnitName  string
	Qty       float32
	UnitPrice float32
	Store     models.StoreMetadata
}

func NewOcrProduct(name string, unitName string, qty float32, unitPrice float32, store models.StoreMetadata) OcrProduct {
	return OcrProduct{
		Name:      name,
		Qty:       qty,
		Price:     qty * unitPrice,
		UnitName:  unitName,
		UnitPrice: unitPrice,
		Store:     store,
	}
}

func (p *OcrProduct) ToDto() dto.OcrProductDto {
	return dto.OcrProductDto{
		ProductName: p.Name,
		UnitPrice:   p.UnitPrice,
		Price:       p.Price,
		Qty:         p.Qty,
		UnitType:    p.UnitName,
		Store:       p.Store.ToDto(),
	}
}

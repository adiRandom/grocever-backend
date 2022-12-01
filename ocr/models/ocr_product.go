package models

import "lib/data/dto"

type OcrProduct struct {
	Name      string
	Price     float32
	UnitName  string
	Qty       float32
	UnitPrice float32
	StoreId   int
}

func NewOcrProduct(name string, unitName string, qty float32, unitPrice float32, storeId int) OcrProduct {
	return OcrProduct{
		Name:      name,
		Qty:       qty,
		Price:     qty * unitPrice,
		UnitName:  unitName,
		UnitPrice: unitPrice,
		StoreId:   storeId,
	}
}

func (p *OcrProduct) ToDto() dto.OcrProductDto {
	return dto.OcrProductDto{
		ProductName:  p.Name,
		ProductPrice: p.Price,
		StoreId:      p.StoreId,
	}
}

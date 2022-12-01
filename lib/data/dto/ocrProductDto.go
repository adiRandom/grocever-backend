package dto

import "fmt"

type OcrProductDto struct {
	ProductName  string  `json:"productName"`
	ProductPrice float32 `json:"productPrice"`
	StoreId      int     `json:"storeId"`
}

func (dto OcrProductDto) String() string {
	return fmt.Sprintf("OcrProductDto: (ProductName: %s ProductPrice: %f StoreId: %d)",
		dto.ProductName,
		dto.ProductPrice,
		dto.StoreId)
}

package dto

type OcrProductDto struct {
	ProductName string        `json:"productName"`
	UnitPrice   float32       `json:"unitPrice"`
	Price       float32       `json:"price"`
	Qty         float32       `json:"qty"`
	UnitType    string        `json:"unitType"`
	Store       StoreMetadata `json:"store"`
}

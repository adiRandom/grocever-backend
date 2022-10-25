package dto

type OcrProductDto struct {
	ProductName  string  `json:"productName"`
	ProductPrice float64 `json:"productPrice"`
	StoreId      int32   `json:"storeId"`
}

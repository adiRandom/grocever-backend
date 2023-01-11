package product_processing

type UserOcrProductDto struct {
	Id        uint    `json:"id"`
	OcrName   string  `json:"ocrName"`
	Qty       float32 `json:"qty"`
	UnitPrice float32 `json:"unitPrice"`
	UnitName  string  `json:"unitName"`
	Price     float32 `json:"price"`
	BestPrice float32 `json:"bestPrice"`
	StoreId   uint    `json:"storeId"`
}

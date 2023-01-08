package product_processing

type UserProductDto struct {
	Id        uint    `json:"id"`
	Name      string  `json:"name"`
	OcrName   string  `json:"ocrName"`
	Qty       float32 `json:"qty"`
	UnitPrice float32 `json:"unitPrice"`
	UnitName  string  `json:"unitName"`
	Price     float32 `json:"price"`
	BestPrice float32 `json:"bestPrice"`
	Url       string  `json:"url"`
	StoreName string  `json:"store"`
	StoreUrl  string  `json:"storeUrl"`
}

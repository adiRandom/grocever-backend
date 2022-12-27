package product_processing

type UserProductDto struct {
	Id        int     `json:"id"`
	Name      string  `json:"name"`
	OcrName   string  `json:"ocrName"`
	Price     float32 `json:"price"`
	BestPrice float32 `json:"bestPrice"`
	Url       string  `json:"url"`
	Store     string  `json:"store"`
}

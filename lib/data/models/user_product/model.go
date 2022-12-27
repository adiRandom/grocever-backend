package user_product

import "lib/data/dto/product_processing"

type Model struct {
	Id        int     `json:"id"`
	Name      string  `json:"name"`
	OcrName   string  `json:"ocrName"`
	Price     float32 `json:"price"`
	BestPrice float32 `json:"bestPrice"`
	Url       string  `json:"url"`
	Store     string  `json:"store"`
	// TODO: Include all data from ocr product
}

func (m *Model) ToDto() product_processing.UserProductDto {
	return product_processing.UserProductDto{
		Id:        m.Id,
		Name:      m.Name,
		OcrName:   m.OcrName,
		Price:     m.Price,
		BestPrice: m.BestPrice,
		Url:       m.Url,
		Store:     m.Store,
	}
}

package entities

type ProductOcrProductSimilarityEntity struct {
	OcrProductName string  `gorm:"column:ocr_product_name;primary_key"`
	ProductId      int     `gorm:"column:product_id;primary_key"`
	Similarity     float64 `gorm:"column:similarity"`
}

func NewProductOcrProductSimilarityEntity(ocrProductName string, productId int, similarity float64) *ProductOcrProductSimilarityEntity {
	return &ProductOcrProductSimilarityEntity{OcrProductName: ocrProductName, ProductId: productId, Similarity: similarity}
}

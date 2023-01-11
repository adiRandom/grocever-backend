package product

type OcrProductModel struct {
	OcrProductName string
	BestPrice      float32
	Products       []*Model
	Related        []*OcrProductModel
}

func NewOcrProductModel(ocrProductName string, bestPrice float32, products []*Model, related []*OcrProductModel) *OcrProductModel {
	return &OcrProductModel{
		OcrProductName: ocrProductName,
		BestPrice:      bestPrice,
		Products:       products,
		Related:        related,
	}
}

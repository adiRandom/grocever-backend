package product

type OcrProductModel struct {
	OcrProductName string
	BestProduct    *Model
	Products       []*Model
	Related        []*OcrProductModel
}

func NewOcrProductModel(ocrProductName string, bestProduct *Model, products []*Model, related []*OcrProductModel) *OcrProductModel {
	return &OcrProductModel{
		OcrProductName: ocrProductName,
		BestProduct:    bestProduct,
		Products:       products,
		Related:        related,
	}
}

package product

type OcrProductModel struct {
	OcrProductName string
	BestPrice      float32
	Products       []Model
	Related        []OcrProductModel
}

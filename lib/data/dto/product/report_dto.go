package product

type ReportDto struct {
	ProductId      uint   `json:"productId"`
	OcrProductName string `json:"ocrProductName"`
	UserId         uint   `json:"userId"`
}

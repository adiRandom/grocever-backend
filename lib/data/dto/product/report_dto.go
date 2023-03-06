package product

type ReportWithUserIdDto struct {
	ProductId      uint   `json:"productId"`
	OcrProductName string `json:"ocrProductName"`
	UserId         uint   `json:"userId"`
}

func NewReportWithUserIdDto(productId uint, ocrProductName string, userId uint) *ReportWithUserIdDto {
	return &ReportWithUserIdDto{ProductId: productId, OcrProductName: ocrProductName, UserId: userId}
}

type ReportDto struct {
	ProductId      uint   `json:"productId"`
	OcrProductName string `json:"ocrProductName"`
}

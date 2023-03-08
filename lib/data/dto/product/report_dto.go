package product

type ReportWithUserIdDto struct {
	ProductId      uint   `json:"productId"`
	OcrProductName string `json:"ocrProductName"`
	UserId         int    `json:"userId"`
}

func NewReportWithUserIdDto(productId uint, ocrProductName string, userId int) *ReportWithUserIdDto {
	return &ReportWithUserIdDto{ProductId: productId, OcrProductName: ocrProductName, UserId: userId}
}

type ReportDto struct {
	ProductId      uint   `json:"productId"`
	OcrProductName string `json:"ocrProductName"`
}

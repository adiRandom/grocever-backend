package product

type ReportWithUserIdDto struct {
	ProductId      int    `json:"productId"`
	OcrProductName string `json:"ocrProductName"`
	UserId         int    `json:"userId"`
}

func NewReportWithUserIdDto(productId int, ocrProductName string, userId int) *ReportWithUserIdDto {
	return &ReportWithUserIdDto{ProductId: productId, OcrProductName: ocrProductName, UserId: userId}
}

type ReportDto struct {
	ProductId      int    `json:"productId"`
	OcrProductName string `json:"ocrProductName"`
}

package ocr

type OcrProductExists struct {
	OcrNames []string `json:"ocrNames" binding:"required"`
}

type OcrProductExistsResponse struct {
	Exists []bool `json:"exists"`
}

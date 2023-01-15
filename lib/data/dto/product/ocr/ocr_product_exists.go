package ocr

type ProductExists struct {
	OcrNames []string `json:"ocrNames" binding:"required"`
}

type ProductExistsResponse struct {
	Exists []bool `json:"exists"`
}

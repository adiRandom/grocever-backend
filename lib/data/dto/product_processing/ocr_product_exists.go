package product_processing

type OcrProductExists struct {
	OcrNames []string `json:"ocrNames" binding:"required"`
}

type OcrProductExistsResponse struct {
	Exists []bool `json:"exists"`
}

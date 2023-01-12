package dto

type CrawlProductDto struct {
	OcrProduct   OcrProductDto    `json:"ocrProduct"`
	CrawlSources []CrawlSourceDto `json:"crawlSources"`
	UserId       int              `json:"userId"`
}

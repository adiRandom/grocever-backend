package dto

import "lib/data/dto/product"

type CrawlProductDto struct {
	OcrProduct   product.UserOcrProductDto `json:"ocrProduct"`
	CrawlSources []CrawlSourceDto          `json:"crawlSources"`
}

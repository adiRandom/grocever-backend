package crawl

import "lib/data/dto/product"

type ProductDto struct {
	OcrProduct   product.UserOcrProductDto `json:"ocrProduct"`
	CrawlSources []SourceDto               `json:"crawlSources"`
}

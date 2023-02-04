package crawl

import "lib/data/dto/product"

type ProductDto struct {
	OcrProduct   product.PurchaseInstalmentDto `json:"ocrProduct"`
	CrawlSources []SourceDto                   `json:"crawlSources"`
}

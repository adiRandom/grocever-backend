package crawl

import "lib/data/dto/product"

type ProductDto struct {
	OcrProduct   product.PurchaseInstalmentWithUserDto `json:"ocrProduct"`
	CrawlSources []SourceDto                           `json:"crawlSources"`
}

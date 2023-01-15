package entities

import (
	"gorm.io/gorm"
	"lib/data/dto/crawl"
	"lib/functional"
)

type ProductRequeueEntity struct {
	gorm.Model
	OcrProductName string
	CrawlSource    crawl.SourceDto `gorm:"embedded"`
}

func NewProductRequeueEntities(product crawl.ProductDto) []ProductRequeueEntity {
	return functional.Map(product.CrawlSources, func(crawlSource crawl.SourceDto) ProductRequeueEntity {
		return ProductRequeueEntity{
			OcrProductName: product.OcrProduct.OcrName,
			CrawlSource:    crawlSource,
		}
	})
}

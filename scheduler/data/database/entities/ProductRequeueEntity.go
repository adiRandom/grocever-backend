package entities

import (
	"gorm.io/gorm"
	"lib/data/dto"
	"lib/functional"
)

type ProductRequeueEntity struct {
	gorm.Model
	OcrProductName string
	CrawlSource    dto.CrawlSourceDto `gorm:"embedded"`
}

func NewProductRequeueEntities(product dto.CrawlProductDto) []ProductRequeueEntity {
	return functional.Map(product.CrawlSources, func(crawlSource dto.CrawlSourceDto) ProductRequeueEntity {
		return ProductRequeueEntity{
			OcrProductName: product.OcrProduct.OcrName,
			CrawlSource:    crawlSource,
		}
	})
}

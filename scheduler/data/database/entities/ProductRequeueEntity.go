package entities

import (
	"gorm.io/gorm"
	"lib/data/dto"
	"lib/functional"
)

type ProductRequeueEntity struct {
	gorm.Model
	Product     dto.OcrProductDto  `gorm:"embedded"`
	CrawlSource dto.CrawlSourceDto `gorm:"embedded"`
}

func NewProductRequeueEntities(product dto.CrawlProductDto) []ProductRequeueEntity {
	return functional.Map(product.CrawlSources, func(crawlSource dto.CrawlSourceDto) ProductRequeueEntity {
		return ProductRequeueEntity{
			Product:     product.OcrProduct,
			CrawlSource: crawlSource,
		}
	})
}

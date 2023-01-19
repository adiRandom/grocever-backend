package entities

import (
	"gorm.io/gorm"
	"lib/data/dto/crawl"
	"lib/functional"
)

type ProductRequeueEntity struct {
	gorm.Model
	OcrProductName string
	CrawlSource    CrawlSource `gorm:"embedded"`
}

func NewProductRequeueEntities(product crawl.ProductDto) []ProductRequeueEntity {
	return functional.Map(product.CrawlSources, func(crawlSource crawl.SourceDto) ProductRequeueEntity {
		return ProductRequeueEntity{
			OcrProductName: product.OcrProduct.OcrName,
			CrawlSource: CrawlSource{
				Url:            crawlSource.Url,
				StoreUrl:       crawlSource.Store.Url,
				StoreId:        crawlSource.Store.StoreId,
				StoreName:      crawlSource.Store.Name,
				OcrHeaderLines: crawlSource.Store.OcrHeaderLines,
			},
		}
	})
}

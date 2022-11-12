package entities

import (
	"gorm.io/gorm"
	dto2 "lib/data/dto"
)

type ProductEntity struct {
	gorm.Model
	Name        string
	CrawlLink   CrawlLinkEntity `gorm:"foreignKey:ProductId"`
	StoreId     int32
	Price       float32
	OcrProducts []*OcrProductEntity `gorm:"many2many:ocr-product_product;"`
}

func NewProductEntities(dto dto2.ProductProcessDto) []ProductEntity {
	productEntities := make([]ProductEntity, len(dto.CrawlResults))
	for i, crawlResult := range dto.CrawlResults {
		productEntities[i] = ProductEntity{
			Name:    crawlResult.ProductName,
			StoreId: crawlResult.StoreId,
			Price:   crawlResult.ProductPrice,
			CrawlLink: CrawlLinkEntity{
				Url:     crawlResult.CrawlUrl,
				StoreId: crawlResult.StoreId,
			},
			OcrProducts: []*OcrProductEntity{&OcrProductEntity{
				OcrProductName: dto.OcrProductDto.ProductName,
			}},
		}
	}

	return productEntities
}

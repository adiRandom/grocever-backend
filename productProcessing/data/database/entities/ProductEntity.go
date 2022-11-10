package entities

import "gorm.io/gorm"

type ProductEntity struct {
	gorm.Model
	Name       string
	CrawlLink  CrawlLinkEntity
	StoreId    int
	Price      float32
	OcrProduct []*OcrProductEntity `gorm:"many2many:ocr-product_product;"`
}

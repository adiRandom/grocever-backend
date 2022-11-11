package entities

import "gorm.io/gorm"

type ProductEntity struct {
	gorm.Model
	Name        string
	CrawlLink   CrawlLinkEntity `gorm:"foreignKey:ProductId"`
	StoreId     int
	Price       float32
	OcrProducts []*OcrProductEntity `gorm:"many2many:ocr-product_product;"`
}

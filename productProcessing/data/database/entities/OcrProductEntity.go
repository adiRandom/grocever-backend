package entities

import "gorm.io/gorm"

// OcrProductEntity Link between products and ocr names
type OcrProductEntity struct {
	gorm.Model
	OcrProductName string
	Product        []*ProductEntity `gorm:"many2many:ocr-product_product;"`
}

package entities

import (
	"gorm.io/gorm"
	"time"
)

// OcrProductEntity Link between products and ocr names
type OcrProductEntity struct {
	OcrProductName string `gorm:"primaryKey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt      `gorm:"index"`
	Products       []*ProductEntity    `gorm:"many2many:ocr-product_product;"`
	Related        []*OcrProductEntity `gorm:"many2many:ocr-product_related;"`
}

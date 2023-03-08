package entities

import (
	"gorm.io/gorm"
)

type MissLink struct {
	gorm.Model
	ProductIdFk      int
	Product          *ProductEntity    `gorm:"foreignKey:ProductIdFk;references:ID"`
	OcrProductNameFk string            `gorm:"size:255"`
	OcrProduct       *OcrProductEntity `gorm:"foreignKey:OcrProductNameFk;references:OcrProductName"`
	UserId           int               `gorm:"unique"`
}

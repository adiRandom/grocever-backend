package entities

import "gorm.io/gorm"

type ProductLinkEntity struct {
	gorm.Model
	ProductID uint
	Url       string
	StoreId   int
}

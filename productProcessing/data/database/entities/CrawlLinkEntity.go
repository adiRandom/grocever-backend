package entities

import "gorm.io/gorm"

type CrawlLinkEntity struct {
	gorm.Model
	Url       string
	StoreId   int32
	ProductId uint
}

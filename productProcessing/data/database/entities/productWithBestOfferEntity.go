package entities

import "gorm.io/gorm"

type Offer struct {
	Price   *float64
	StoreId *int
}

type ProductWithBestOfferEntity struct {
	gorm.Model
	Name  string `gorm:"not null;index"`
	Offer Offer  `gorm:"embedded;embeddedPrefix:offer_"`
}

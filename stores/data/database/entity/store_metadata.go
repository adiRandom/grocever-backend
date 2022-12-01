package entity

import "gorm.io/gorm"

type StoreMetadata struct {
	gorm.Model
	StoreId        int
	Name           string
	OcrHeaderLines int
}

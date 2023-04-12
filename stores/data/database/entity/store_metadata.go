package entity

import (
	"gorm.io/gorm"
	"lib/data/models"
)

type StoreMetadata struct {
	gorm.Model
	StoreId        int
	Name           string
	OcrHeaderLines int
	Url            string
}

func (s *StoreMetadata) ToModel() models.StoreMetadata {
	return models.StoreMetadata{
		StoreId: s.StoreId,
		Name:    s.Name,
		Url:     s.Url,
	}
}

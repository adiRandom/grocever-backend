package entities

import (
	"gorm.io/gorm"
	"lib/data/dto"
)

type CrawlLinkEntity struct {
	gorm.Model
	Url       string
	StoreId   int32
	ProductId uint
}

func (entity CrawlLinkEntity) ToDto() dto.CrawlSourceDto {
	return dto.CrawlSourceDto{
		Url:     entity.Url,
		StoreId: int(entity.StoreId),
	}
}

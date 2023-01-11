package entities

import (
	"gorm.io/gorm"
	"lib/data/dto"
	"lib/data/models/crawl"
)

type CrawlLinkEntity struct {
	gorm.Model
	Url       string
	StoreId   int
	ProductId uint
}

func (entity CrawlLinkEntity) ToDto() dto.CrawlSourceDto {
	return dto.CrawlSourceDto{
		Url:     entity.Url,
		StoreId: int(entity.StoreId),
	}
}

func (entity CrawlLinkEntity) ToModel() crawl.LinkModel {
	return crawl.LinkModel{
		Id:        entity.ID,
		Url:       entity.Url,
		StoreId:   entity.StoreId,
		ProductId: entity.ProductId,
	}
}

func NewCrawlLinkEntityFromModel(model crawl.LinkModel) CrawlLinkEntity {
	entity := CrawlLinkEntity{
		Url:       model.Url,
		StoreId:   model.StoreId,
		ProductId: model.ProductId,
	}

	if model.Id != -1 {
		entity.ID = uint(model.Id)
	}
	return entity
}

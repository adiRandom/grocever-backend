package entities

import (
	"gorm.io/gorm"
	"lib/data/dto"
	"lib/data/models/crawl"
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

func (entity CrawlLinkEntity) ToModel() crawl.LinkModel {
	return crawl.LinkModel{
		Id:        entity.ID,
		Url:       entity.Url,
		StoreId:   entity.StoreId,
		ProductId: entity.ProductId,
	}
}

func NewCrawlLinkEntityFromModel(model crawl.LinkModel) CrawlLinkEntity {
	return CrawlLinkEntity{
		Model:     gorm.Model{ID: model.Id},
		Url:       model.Url,
		StoreId:   model.StoreId,
		ProductId: model.ProductId,
	}
}

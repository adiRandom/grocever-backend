package entities

import (
	"gorm.io/gorm"
	"lib/data/models/crawl"
)

type CrawlLinkEntity struct {
	gorm.Model
	Url       string
	StoreId   int
	ProductId uint
}

func (entity CrawlLinkEntity) ToModel() crawl.LinkModel {
	return crawl.LinkModel{
		Id:        int(entity.ID),
		Url:       entity.Url,
		StoreId:   entity.StoreId,
		ProductId: int(entity.ProductId),
	}
}

func NewCrawlLinkEntityFromModel(model crawl.LinkModel) *CrawlLinkEntity {

	if model.ProductId == -1 {
		return nil
	}

	entity := CrawlLinkEntity{
		Url:       model.Url,
		StoreId:   model.StoreId,
		ProductId: uint(model.ProductId),
	}

	if model.Id != -1 {
		entity.ID = uint(model.Id)
	}
	return &entity
}

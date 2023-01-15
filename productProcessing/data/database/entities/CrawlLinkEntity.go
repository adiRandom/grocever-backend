package entities

import (
	"gorm.io/gorm"
	"lib/data/models/crawl"
	"lib/helpers"
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

func NewCrawlLinkEntityFromModel(model crawl.LinkModel) (*CrawlLinkEntity, error) {

	if model.ProductId == -1 {
		return nil, helpers.Error{Msg: "ProductId cannot be -1"}
	}

	entity := CrawlLinkEntity{
		Url:       model.Url,
		StoreId:   model.StoreId,
		ProductId: uint(model.ProductId),
	}

	if model.Id != -1 {
		entity.ID = uint(model.Id)
	}
	return &entity, nil
}

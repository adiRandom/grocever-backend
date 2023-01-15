package entities

import (
	"gorm.io/gorm"
	"lib/data/models/product"
	"lib/functional"
)

type ProductEntity struct {
	gorm.Model
	Name        string
	CrawlLink   CrawlLinkEntity `gorm:"foreignKey:ProductId"`
	StoreId     int
	Price       float32
	UnityType   string
	OcrProducts []*OcrProductEntity `gorm:"many2many:ocr-product_product;"`
}

func (entity ProductEntity) ToModel() product.Model {
	return product.Model{
		ID:        int(entity.ID),
		Name:      entity.Name,
		StoreId:   entity.StoreId,
		Price:     entity.Price,
		UnityType: entity.UnityType,
		CrawlLink: entity.CrawlLink.ToModel(),
		OcrProducts: functional.Map(entity.OcrProducts,
			func(ocrProductEntity *OcrProductEntity) *product.OcrProductModel {
				model := ocrProductEntity.ToModel(false, false)
				return &model
			},
		),
	}
}

func NewProductEntityFromModel(model product.Model) (*ProductEntity, error) {
	crawlLinkEntity, err := NewCrawlLinkEntityFromModel(model.CrawlLink)
	if err != nil {
		return nil, err
	}

	entity := ProductEntity{
		Name:      model.Name,
		StoreId:   model.StoreId,
		Price:     model.Price,
		UnityType: model.UnityType,
		CrawlLink: *crawlLinkEntity,
	}

	if model.ID != -1 {
		entity.ID = uint(model.ID)
	}
	return &entity, nil
}

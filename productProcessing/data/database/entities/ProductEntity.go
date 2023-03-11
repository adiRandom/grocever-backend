package entities

import (
	"gorm.io/gorm"
	"lib/data/models/product"
	"lib/functional"
)

type ProductEntity struct {
	gorm.Model
	Name        string
	ImageUrl    string
	CrawlLink   *CrawlLinkEntity `gorm:"foreignKey:ProductId"`
	StoreId     int
	Price       float32
	UnityType   string
	OcrProducts []*OcrProductEntity `gorm:"many2many:ocr-product_product;"`
}

func (entity ProductEntity) ToModel() product.Model {
	var ocrProducts []*product.OcrProductModel = nil
	if entity.OcrProducts != nil {
		ocrProducts = functional.Map(entity.OcrProducts, func(ocrProduct *OcrProductEntity) *product.OcrProductModel {
			model := ocrProduct.ToModel(false, false)
			return &model
		})
	}
	return product.Model{
		ID:          int(entity.ID),
		ImageUrl:    entity.ImageUrl,
		Name:        entity.Name,
		StoreId:     entity.StoreId,
		Price:       entity.Price,
		UnityType:   entity.UnityType,
		CrawlLink:   entity.CrawlLink.ToModel(),
		OcrProducts: ocrProducts,
	}
}

func NewProductEntityFromModel(model product.Model) *ProductEntity {
	crawlLinkEntity := NewCrawlLinkEntityFromModel(model.CrawlLink)

	entity := ProductEntity{
		Name:      model.Name,
		StoreId:   model.StoreId,
		Price:     model.Price,
		UnityType: model.UnityType,
		CrawlLink: crawlLinkEntity,
		ImageUrl:  model.ImageUrl,
	}

	if model.ID != -1 {
		entity.ID = uint(model.ID)
	}
	return &entity
}

package entities

import (
	"gorm.io/gorm"
	dto2 "lib/data/dto"
	"lib/data/models/product"
	"lib/functional"
)

type ProductEntity struct {
	gorm.Model
	Name        string
	CrawlLink   CrawlLinkEntity `gorm:"foreignKey:ProductId"`
	StoreId     int32
	Price       float32
	UnityType   string
	OcrProducts []*OcrProductEntity `gorm:"many2many:ocr-product_product;"`
}

func NewProductEntities(dto dto2.ProductProcessDto) []ProductEntity {
	productEntities := make([]ProductEntity, len(dto.CrawlResults))
	for i, crawlResult := range dto.CrawlResults {
		productEntities[i] = ProductEntity{
			Name:    crawlResult.ProductName,
			StoreId: crawlResult.StoreId,
			Price:   crawlResult.ProductPrice,
			CrawlLink: CrawlLinkEntity{
				Url:     crawlResult.CrawlUrl,
				StoreId: crawlResult.StoreId,
			},
			OcrProducts: []*OcrProductEntity{&OcrProductEntity{
				OcrProductName: dto.OcrProductDto.ProductName,
			}},
		}
	}

	return productEntities
}

func (entity ProductEntity) ToModel() product.Model {
	return product.Model{
		ID:        uint(entity.ID),
		Name:      entity.Name,
		StoreId:   entity.StoreId,
		Price:     entity.Price,
		UnityType: entity.UnityType,
		CrawlLink: entity.CrawlLink.ToModel(),
		OcrProducts: functional.Map(entity.OcrProducts,
			func(ocrProductEntity *OcrProductEntity) product.OcrProductModel {
				return ocrProductEntity.ToModel(false, false)
			},
		),
	}
}

func NewProductEntityFromModel(model product.Model) ProductEntity {
	return ProductEntity{
		Model:     gorm.Model{ID: uint(model.ID)},
		Name:      model.Name,
		StoreId:   model.StoreId,
		Price:     model.Price,
		UnityType: model.UnityType,
		CrawlLink: NewCrawlLinkEntityFromModel(model.CrawlLink),
	}
}

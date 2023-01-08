package product

import (
	"lib/data/dto/product_processing"
	"lib/data/models"
)

type UserOcrProductModel struct {
	Id         uint
	Qty        float32
	Price      float32
	UserId     string
	OcrProduct OcrProductModel
	Product    Model
	Store      models.StoreMetadata
}

func (m *UserOcrProductModel) ToDto() product_processing.UserProductDto {
	return product_processing.UserProductDto{
		Id:        m.Id,
		Name:      m.Product.Name,
		OcrName:   m.OcrProduct.OcrProductName,
		Qty:       m.Qty,
		UnitPrice: m.Product.Price,
		UnitName:  m.Product.UnityType,
		Price:     m.Price,
		BestPrice: m.OcrProduct.BestPrice,
		Url:       m.Product.CrawlLink.Url,
		StoreName: m.Store.Name,
		StoreUrl:  m.Store.Url,
	}
}

package models

import (
	"lib/data/dto/product"
	models "lib/data/models/product"
	"lib/functional"
)

type UserProduct struct {
	Id                  uint                             `json:"id"`
	Name                string                           `json:"name"`
	BestPrice           float32                          `json:"bestPrice"`
	PurchaseInstalments []models.PurchaseInstalmentModel `json:"purchaseInstalments"`
	BestStoreId         uint                             `json:"bestStoreId"`
	BestStoreName       string                           `json:"bestStoreName"`
	BestStoreUrl        string                           `json:"bestStoreUrl"`
	BestProductUrl      string                           `json:"bestProductUrl"`
	ImageUrl            string                           `json:"imageUrl"`
}

func NewUserProduct(
	id uint,
	name string,
	bestPrice float32,
	purchaseInstalments []models.PurchaseInstalmentModel,
	bestStoreId uint,
	bestStoreName string,
	bestStoreUrl string,
	bestProductUrl string,
	imageUrl string,
) *UserProduct {
	return &UserProduct{
		Id:                  id,
		Name:                name,
		BestPrice:           bestPrice,
		PurchaseInstalments: purchaseInstalments,
		BestStoreId:         bestStoreId,
		BestStoreName:       bestStoreName,
		BestStoreUrl:        bestStoreUrl,
		BestProductUrl:      bestProductUrl,
		ImageUrl:            imageUrl,
	}
}

func (p *UserProduct) ToDto() product.UserProductDto {
	return product.UserProductDto{
		Id:        p.Id,
		Name:      p.Name,
		BestPrice: p.BestPrice,
		PurchaseInstalments: functional.Map(p.PurchaseInstalments, func(p models.PurchaseInstalmentModel) product.PurchaseInstalmentDto {
			return p.ToDto()
		}),
		BestStoreId:    p.BestStoreId,
		BestStoreName:  p.BestStoreName,
		BestStoreUrl:   p.BestStoreUrl,
		BestProductUrl: p.BestProductUrl,
		ImageUrl:       p.ImageUrl,
	}
}

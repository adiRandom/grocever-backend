package entities

import (
	"gorm.io/gorm"
	"lib/data/dto"
)

type ProductRequeueEntity struct {
	gorm.Model
	Product dto.CrawlProductDto
}

func NewProductRequeueEntity(product dto.CrawlProductDto) *ProductRequeueEntity {
	return &ProductRequeueEntity{
		Product: product,
	}
}

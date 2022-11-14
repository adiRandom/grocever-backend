package entities

import (
	"gorm.io/gorm"
	"lib/data/models"
)

type ProductRequeueEntity struct {
	gorm.Model
	Product models.CrawlerResult
}

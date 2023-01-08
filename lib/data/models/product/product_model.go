package product

import (
	"lib/data/models/crawl"
)

type Model struct {
	ID          uint
	Name        string
	CrawlLink   crawl.LinkModel
	StoreId     int32
	Price       float32
	UnityType   string
	OcrProducts []OcrProductModel
}

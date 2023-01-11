package crawl

import (
	"lib/data/dto"
)

type CrawlerResult struct {
	ProductName  string
	ProductPrice float32
	Store        dto.StoreMetadata
	CrawlUrl     string
	// ImageUrl string
}

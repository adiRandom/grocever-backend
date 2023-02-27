package crawl

import (
	"lib/data/dto/store"
)

type ResultDto struct {
	ProductName  string
	ProductPrice float32
	Store        store.MetadataDto
	CrawlUrl     string
	ImageUrl     string
}

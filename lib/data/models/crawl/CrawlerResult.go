package crawl

import (
	"lib/data/dto/crawl"
	"lib/data/models"
)

type ResultModel struct {
	ProductName  string
	ProductPrice float32
	Store        models.StoreMetadata
	CrawlUrl     string
	ImageUrl     string
}

func (c *ResultModel) ToDto() crawl.ResultDto {
	return crawl.ResultDto{
		ProductName:  c.ProductName,
		ProductPrice: c.ProductPrice,
		Store:        c.Store.ToDto(),
		CrawlUrl:     c.CrawlUrl,
		ImageUrl:     c.ImageUrl,
	}
}

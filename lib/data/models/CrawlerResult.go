package models

import "fmt"

type CrawlerResult struct {
	ProductName  string
	ProductPrice float32
	StoreId      int32
	CrawlUrl     string
	// TODO: Add
	// UnitPrice float64
	// Unit string
	// ImageUrl string
}

func (res CrawlerResult) String() string {
	return fmt.Sprintf("Product: %s at price: %f from store %d", res.ProductName, res.ProductPrice, res.StoreId)
}

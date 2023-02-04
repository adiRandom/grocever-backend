package dto

import (
	"lib/data/dto/crawl"
	"lib/data/dto/product"
)

type ProductProcessDto struct {
	OcrProduct   product.PurchaseInstalmentWithUserDto `json:"ocrProduct"`
	CrawlResults []crawl.ResultDto                     `json:"crawlResult"`
}

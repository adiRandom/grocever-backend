package dto

import (
	"lib/data/dto/product"
	"lib/data/models/crawl"
)

type ProductProcessDto struct {
	OcrProduct   product.UserOcrProductDto `json:"ocrProduct"`
	CrawlResults []crawl.CrawlerResult     `json:"crawlResult"`
}

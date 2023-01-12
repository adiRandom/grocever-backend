package dto

import (
	"lib/data/models/crawl"
)

type ProductProcessDto struct {
	OcrProductDto OcrProductDto         `json:"ocrProduct"`
	CrawlResults  []crawl.CrawlerResult `json:"crawlResult"`
	UserId        int                   `json:"userId"`
}

package services

import (
	"lib/data/models/crawl"
)

type Crawler interface {
	ScrapeProductPage(url string, resCh chan crawl.ResultModel)
}

package services

import (
	"lib/data/models/crawl"
)

type Crawler interface {
	// ScrapeProductPage
	// Send to channel a struct with empty string for crawl url if error occurred
	ScrapeProductPage(url string, resCh chan crawl.ResultModel)
}

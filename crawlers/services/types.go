package services

import (
	"lib/data/models"
)

type Crawler interface {
	ScrapeProductPage(url string, resCh chan models.CrawlerResult)
}

package services

import "crawlers/models"

type Crawler interface {
	ScrapeProductPage(url string, resCh chan models.CrawlerResult)
}

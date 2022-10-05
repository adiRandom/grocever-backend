package crawlers

import "dealScraper/crawlers/models"

type Crawler interface {
	ScrapeProductPage(url string, resCh chan models.CrawlerResult)
}

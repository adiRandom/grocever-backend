package main

import (
	"dealScraper/crawlers/models"
	"dealScraper/crawlers/services"
	"fmt"
)

func main() {
	resChan := make(chan models.CrawlerResult)
	go crawlers.ScrapeProductPage("https://www.auchan.ro/store/Lapte-integral-Zuzu%2C-3-5-grasime%2C-1L/p/000888", resChan)
	res := <-resChan
	fmt.Println(res.String())
}

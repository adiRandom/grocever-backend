package main

import (
	"dealScraper/crawlers/models"
	crawlers "dealScraper/crawlers/services"
	"fmt"
)

func main() {
	resChan := make(chan models.CrawlerResult)
	go crawlers.CoraCrawler{}.ScrapeProductPage("https://www.cora.ro/zuzu-lapte-de-consum-3-5-grasime-1-l-2077870.html", resChan)
	res := <-resChan
	fmt.Println(res.String())
}

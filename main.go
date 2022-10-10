package main

import (
	"dealScraper/crawlers/models"
	crawlers "dealScraper/crawlers/services"
	"fmt"
)

func main() {
	resChan := make(chan models.CrawlerResult)
	go crawlers.MegaImageCrawler{}.ScrapeProductPage("https://www.mega-image.ro/ro-ro/Lactate-si-oua/Lapte/Lapte-de-consum-semidegresat/Lapte-de-consum-1-5-grasime-1L/p/36688", resChan)
	res := <-resChan
	fmt.Println(res.String())
}

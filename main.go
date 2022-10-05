package main

import (
	"dealScraper/crawlers/models"
	crawlers "dealScraper/crawlers/services"
	"fmt"
)

func main() {
	resChan := make(chan models.CrawlerResult)
	go crawlers.FreshfulCrawler{}.ScrapeProductPage("https://www.freshful.ro/p/100004014-zuzu-lapte-1-5-grasime-1l?gclid=Cj0KCQjw1vSZBhDuARIsAKZlijQ-cBKNqX91OI-3D2OPhS1fl3UpfBf9dahoUeGbcVil16e-EHJ00QsaAjNGEALw_wcB", resChan)
	res := <-resChan
	fmt.Println(res.String())
}

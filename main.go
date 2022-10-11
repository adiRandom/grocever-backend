package main

import (
	"dealScraper/search/services"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	res, err := services.GoogleSearchService{}.SearchCrawlSources("lapte zuzu")
	if err != nil {
		panic(err)
	}

	for _, item := range res {
		println(item.Url + " " + string(item.StoreId))
	}
}

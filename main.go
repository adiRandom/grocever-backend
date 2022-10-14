package main

import (
	"dealScraper/data/database"
	"dealScraper/data/database/entities"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	err = database.InitDatabase(&entities.ProductWithBestOfferEntity{})
	if err != nil {
		panic(err)
	}
}

package main

import (
	"github.com/joho/godotenv"
	"productProcessing/data/database"
	"productProcessing/data/database/entities"
	"productProcessing/data/database/repositories"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	database.InitDatabase(&entities.ProductEntity{}, &entities.CrawlLinkEntity{}, &entities.OcrProductEntity{})
	repositories.GetProductRepository().Save(
		entities.ProductEntity{
			Name:    "test",
			StoreId: 1,
			Price:   1,
			CrawlLink: entities.CrawlLinkEntity{
				Url:     "test",
				StoreId: 1,
			},
			OcrProduct: make([]*entities.OcrProductEntity, 0),
		},
	)

	x, err := repositories.GetProductRepository().GetAllWithCrawlLink()
	if err != nil {
		panic(err)
	}
	println(x[0].CrawlLink.Url)
}

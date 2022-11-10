package main

import (
	"productProcessing/data/database"
	"productProcessing/data/database/entities"
	"productProcessing/data/database/repositories"
)

func main() {
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

	x, err := repositories.GetProductRepository().GetAll()
	if err != nil {
		panic(err)
	}
	println(x)
}

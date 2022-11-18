package test

import (
	"productProcessing/data/database/entities"
	"productProcessing/data/database/repositories"
)

func InitDbTestData() {
	ocrProductEntities := [1]*entities.OcrProductEntity{
		&entities.OcrProductEntity{
			OcrProductName: "test",
		},
	}
	repositories.GetProductRepository().Create(
		&entities.ProductEntity{
			Name:    "test",
			StoreId: 1,
			Price:   1.5,
			CrawlLink: entities.CrawlLinkEntity{
				Url:     "test",
				StoreId: 1,
			},
			OcrProducts: ocrProductEntities[:],
		},
	)

	repositories.GetProductRepository().Create(
		&entities.ProductEntity{
			Name:    "test2",
			StoreId: 2,
			Price:   2.7,
			CrawlLink: entities.CrawlLinkEntity{
				Url:     "test",
				StoreId: 2,
			},
			OcrProducts: ocrProductEntities[:],
		},
	)

	ocrProductEntities = [1]*entities.OcrProductEntity{
		&entities.OcrProductEntity{
			OcrProductName: "tes",
		},
	}

	repositories.GetProductRepository().Create(
		&entities.ProductEntity{
			Name:    "test",
			StoreId: 1,
			Price:   1.1,
			CrawlLink: entities.CrawlLinkEntity{
				Url:     "test",
				StoreId: 1,
			},
			OcrProducts: ocrProductEntities[:],
		},
	)
}

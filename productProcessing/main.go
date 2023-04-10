package main

import (
	"context"
	"github.com/joho/godotenv"
	"lib/data/database"
	"os"
	"productProcessing/data/database/entities"
	"productProcessing/gateways/api"
	"productProcessing/gateways/events"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	err = database.InitDatabase(
		&entities.ProductEntity{},
		&entities.CrawlLinkEntity{},
		&entities.OcrProductEntity{},
		&entities.PurchaseInstalment{},
		&entities.MissLink{},
		&entities.ProductOcrProductSimilarityEntity{},
	)
	if err != nil {
		return
	}

	println("Started")
	go events.GetRabbitMqBroker().Start(context.Background())

	r := api.GetRouter()
	r.Run(os.Getenv("API_PORT"))

	// Testing

	//service := services.NewProductService(
	//	repositories.GetProductRepository(
	//		repositories.GetMissLinkRepository(),
	//		repositories.GetOcrProductRepository(repositories.GetMissLinkRepository()),
	//	),
	//	repositories.GetOcrProductRepository(repositories.GetMissLinkRepository()),
	//	repositories.GetUserProductRepository(),
	//)
	//
	//dto1 := dto.ProductProcessDto{
	//	OcrProduct: product.PurchaseInstalmentWithUserDto{
	//		PurchaseInstalmentDto: product.PurchaseInstalmentDto{
	//			Id:        -1,
	//			OcrName:   "1",
	//			Price:     100,
	//			Qty:       1,
	//			UnitPrice: 100,
	//			UnitName:  "kg",
	//			Store: store.MetadataDto{
	//				StoreId: 1,
	//				Name:    "test",
	//			},
	//		},
	//		UserId: 4,
	//	},
	//	CrawlResults: []crawl.ResultDto{
	//		{
	//			ProductName:  "A",
	//			ProductPrice: 3,
	//			Store: store.MetadataDto{
	//				StoreId: 1,
	//				Name:    "test",
	//			},
	//		},
	//	},
	//}
	//
	//dto2 := dto.ProductProcessDto{
	//	OcrProduct: product.PurchaseInstalmentWithUserDto{
	//		PurchaseInstalmentDto: product.PurchaseInstalmentDto{
	//			Id:        -1,
	//			OcrName:   "2",
	//			Price:     100,
	//			Qty:       1,
	//			UnitPrice: 100,
	//			UnitName:  "kg",
	//			Store: store.MetadataDto{
	//				StoreId: 1,
	//				Name:    "test",
	//			},
	//		},
	//		UserId: 4,
	//	},
	//	CrawlResults: []crawl.ResultDto{
	//		crawl.ResultDto{
	//			ProductName:  "B",
	//			ProductPrice: 1,
	//			Store: store.MetadataDto{
	//				StoreId: 1,
	//				Name:    "test",
	//			},
	//		},
	//	},
	//}
	//
	//dto3 := dto.ProductProcessDto{
	//	OcrProduct: product.PurchaseInstalmentWithUserDto{
	//		PurchaseInstalmentDto: product.PurchaseInstalmentDto{
	//			Id:        -1,
	//			OcrName:   "3",
	//			Price:     100,
	//			Qty:       1,
	//			UnitPrice: 100,
	//			UnitName:  "kg",
	//			Store: store.MetadataDto{
	//				StoreId: 1,
	//				Name:    "test",
	//			},
	//		},
	//		UserId: 4,
	//	},
	//	CrawlResults: []crawl.ResultDto{
	//		crawl.ResultDto{
	//			ProductName:  "C",
	//			ProductPrice: 2,
	//			Store: store.MetadataDto{
	//				StoreId: 1,
	//				Name:    "test",
	//			},
	//		},
	//	},
	//}
	//
	//dto4 := dto.ProductProcessDto{
	//	OcrProduct: product.PurchaseInstalmentWithUserDto{
	//		PurchaseInstalmentDto: product.PurchaseInstalmentDto{
	//			Id:        -1,
	//			OcrName:   "2",
	//			Price:     100,
	//			Qty:       1,
	//			UnitPrice: 100,
	//			UnitName:  "kg",
	//			Store: store.MetadataDto{
	//				StoreId: 1,
	//				Name:    "test",
	//			},
	//		},
	//		UserId: 4,
	//	},
	//	CrawlResults: []crawl.ResultDto{
	//		crawl.ResultDto{
	//			ProductName:  "C",
	//			ProductPrice: 2,
	//			Store: store.MetadataDto{
	//				StoreId: 1,
	//				Name:    "test",
	//			},
	//		},
	//	},
	//}
	//
	//dto5 := dto.ProductProcessDto{
	//	OcrProduct: product.PurchaseInstalmentWithUserDto{
	//		PurchaseInstalmentDto: product.PurchaseInstalmentDto{
	//			Id:        -1,
	//			OcrName:   "1",
	//			Price:     100,
	//			Qty:       1,
	//			UnitPrice: 100,
	//			UnitName:  "kg",
	//			Store: store.MetadataDto{
	//				StoreId: 1,
	//				Name:    "test",
	//			},
	//		},
	//		UserId: 4,
	//	},
	//	CrawlResults: []crawl.ResultDto{
	//		crawl.ResultDto{
	//			ProductName:  "B",
	//			ProductPrice: 1,
	//			Store: store.MetadataDto{
	//				StoreId: 1,
	//				Name:    "test",
	//			},
	//		},
	//	},
	//}
	//
	//dto6 := dto.ProductProcessDto{
	//	OcrProduct: product.PurchaseInstalmentWithUserDto{
	//		PurchaseInstalmentDto: product.PurchaseInstalmentDto{
	//			Id:        -1,
	//			OcrName:   "1",
	//			Price:     100,
	//			Qty:       1,
	//			UnitPrice: 100,
	//			UnitName:  "kg",
	//			Store: store.MetadataDto{
	//				StoreId: 1,
	//				Name:    "test",
	//			},
	//		},
	//		UserId: 4,
	//	},
	//	CrawlResults: []crawl.ResultDto{
	//		crawl.ResultDto{
	//			ProductName:  "D",
	//			ProductPrice: 4,
	//			Store: store.MetadataDto{
	//				StoreId: 1,
	//				Name:    "test",
	//			},
	//		},
	//	},
	//}
	//
	//dto7 := dto.ProductProcessDto{
	//	OcrProduct: product.PurchaseInstalmentWithUserDto{
	//		PurchaseInstalmentDto: product.PurchaseInstalmentDto{
	//			Id:        -1,
	//			OcrName:   "2",
	//			Price:     100,
	//			Qty:       1,
	//			UnitPrice: 100,
	//			UnitName:  "kg",
	//			Store: store.MetadataDto{
	//				StoreId: 1,
	//				Name:    "test",
	//			},
	//		},
	//		UserId: 4,
	//	},
	//	CrawlResults: []crawl.ResultDto{
	//		crawl.ResultDto{
	//			ProductName:  "D",
	//			ProductPrice: 4,
	//			Store: store.MetadataDto{
	//				StoreId: 1,
	//				Name:    "test",
	//			},
	//		},
	//	},
	//}
	//
	//service.ProcessCrawlProduct(dto1)
	//service.ProcessCrawlProduct(dto2)
	//service.ProcessCrawlProduct(dto3)
	//service.ProcessCrawlProduct(dto4)
	//service.ProcessCrawlProduct(dto5)
	//service.ProcessCrawlProduct(dto6)
	//service.ProcessCrawlProduct(dto7)
	//
	//r := repositories.GetMissLinkRepository()
	//r.CreateMany([]entities.MissLink{
	//	entities.MissLink{
	//		UserId:           4,
	//		ProductIdFk:      2,
	//		OcrProductNameFk: "1",
	//	},
	//	entities.MissLink{
	//		UserId:           1,
	//		ProductIdFk:      2,
	//		OcrProductNameFk: "1",
	//	},
	//	entities.MissLink{
	//		UserId:           2,
	//		ProductIdFk:      2,
	//		OcrProductNameFk: "1",
	//	},
	//})
	//
	//err = repositories.GetProductRepository(
	//	repositories.GetMissLinkRepository(),
	//	repositories.GetOcrProductRepository(repositories.GetMissLinkRepository()),
	//).BreakProductLinkAsync(2, "1")
	//
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//// Wait for async process to finish
	//infiniteLoop := make(chan bool)
	//<-infiniteLoop
}

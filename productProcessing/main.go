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

	err = database.InitDatabase(&entities.ProductEntity{}, &entities.CrawlLinkEntity{}, &entities.OcrProductEntity{}, &entities.PurchaseInstalment{}, &entities.MissLink{})
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
	//		repositories.GetOcrProductRepository(),
	//	),
	//	repositories.GetOcrProductRepository(),
	//	repositories.GetUserProductRepository(),
	//)
	//
	//dto1 := dto.ProductProcessDto{
	//	OcrProduct: product.PurchaseInstalmentWithUserDto{
	//		PurchaseInstalmentDto: product.PurchaseInstalmentDto{
	//			Id:        -1,
	//			OcrName:   "test",
	//			Price:     100,
	//			Qty:       1,
	//			UnitPrice: 1,
	//			UnitName:  "kg",
	//			Store: store.MetadataDto{
	//				StoreId: 1,
	//				Name:    "test",
	//			},
	//		},
	//		UserId: 1,
	//	},
	//	// TODO: Write dtos to test best price updates afrer link and after unlink
	//	CrawlResults: []crawl.ResultDto{
	//		crawl.ResultDto{
	//			ProductName:  "test",
	//			ProductPrice: 100,
	//			Store: store.MetadataDto{
	//				StoreId: 1,
	//				Name:    "test",
	//			},
	//		},
	//		crawl.ResultDto{
	//			ProductName:  "test2",
	//			ProductPrice: 75,
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
	//			OcrName:   "test2",
	//			Price:     100,
	//			Qty:       1,
	//			UnitPrice: 1,
	//			UnitName:  "kg",
	//			Store: store.MetadataDto{
	//				StoreId: 1,
	//				Name:    "test",
	//			},
	//		},
	//		UserId: 1,
	//	},
	//	// TODO: Write dtos to test best price updates afrer link and after unlink
	//	CrawlResults: []crawl.ResultDto{
	//		crawl.ResultDto{
	//			ProductName:  "test",
	//			ProductPrice: 100,
	//			Store: store.MetadataDto{
	//				StoreId: 1,
	//				Name:    "test",
	//			},
	//		},
	//		crawl.ResultDto{
	//			ProductName:  "test3",
	//			ProductPrice: 75,
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

	//repositories.GetProductRepository(
	//	repositories.GetMissLinkRepository(),
	//	repositories.GetOcrProductRepository(),
	//).BreakProductLink(1, "test2")
}

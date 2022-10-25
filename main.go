package main

import (
	"context"
	"dealScraper/lib/data/dto"
	"dealScraper/ocr"
	"dealScraper/search/services"
)

func main() {
	//err := godotenv.Load()
	//if err != nil {
	//	panic(err)
	//}
	//err = database.InitDatabase(&entities.ProductWithBestOfferEntity{})
	//if err != nil {
	//	panic(err)
	//}

	go ocr.SendOcrProductToQueue()
	services.ListenForSearchRequests(context.Background(), func(ocrProduct dto.OcrProductDto) {
		println(ocrProduct.ProductName)
	})
}

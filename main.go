package main

import (
	"context"
	"dealScraper/crawlers/messages"
	"dealScraper/crawlers/test"
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

	test.ProduceCrawlMessages()
	ctx := context.Background()
	go messages.GetRabbitMqBroker().ListenAndHandleRequests(ctx)

	select {}
}

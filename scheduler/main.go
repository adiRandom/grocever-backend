package main

import (
	"context"
	"github.com/joho/godotenv"
	"lib/data/database"
	"scheduler/data/database/entities"
	"scheduler/gateways/events"
	"scheduler/services/crawl"
	"scheduler/test"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	err = database.InitDatabase(&entities.ProductRequeueEntity{})
	if err != nil {
		panic(err)
	}

	crawlScheduler := crawl.GetScheduler()
	defer crawlScheduler.Close()

	broker := events.GetRabbitMqBroker()
	go broker.Start(context.Background())

	service := crawl.GetRequeueService()
	service.StartCronRequeue()

	test.ProduceRequeueMessages()

	select {}
}

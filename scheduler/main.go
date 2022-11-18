package main

import (
	"context"
	"github.com/joho/godotenv"
	"scheduler/gateways/events"
	"scheduler/services/crawl"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	crawlScheduler := crawl.GetScheduler()
	defer crawlScheduler.Close()

	broker := events.GetRabbitMqBroker()
	go broker.Start(context.Background())

	service := crawl.GetRequeueService()
	service.StartCronRequeue()

	select {}
}

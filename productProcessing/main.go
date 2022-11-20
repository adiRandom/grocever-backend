package main

import (
	"context"
	"github.com/joho/godotenv"
	"lib/data/database"
	"productProcessing/data/database/entities"
	"productProcessing/gateways/events"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	err = database.InitDatabase(&entities.ProductEntity{}, &entities.CrawlLinkEntity{}, &entities.OcrProductEntity{})
	if err != nil {
		return
	}

	println("Started")
	events.GetRabbitMqBroker().Start(context.Background())
}

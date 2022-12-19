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

	err = database.InitDatabase(&entities.ProductEntity{}, &entities.CrawlLinkEntity{}, &entities.OcrProductEntity{})
	if err != nil {
		return
	}

	println("Started")
	go events.GetRabbitMqBroker().Start(context.Background())

	r := api.GetRouter()
	r.Run(os.Getenv("API_PORT"))
}

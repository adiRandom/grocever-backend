package main

import (
	"context"
	"github.com/joho/godotenv"
	"ocr/gateways/api"
	"ocr/gateways/events"
	mock "ocr/test/events"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	println(os.Getenv("API_PORT"))

	mock.PublishMockProducts()

	router := api.GetRouter()
	broker := events.GetRabbitMqBroker()
	ctx := context.Background()
	go broker.Start(ctx)
	router.Run(os.Getenv("API_PORT"))
}

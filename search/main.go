package main

import (
	"context"
	"github.com/joho/godotenv"
	"search/events"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	events.GetRabbitMqBroker().Start(context.Background())
}

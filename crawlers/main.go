package main

import (
	"context"
	"crawlers/gateways/events"
)

func main() {
	broker := events.GetRabbitMqBroker()
	broker.Start(context.Background())
}

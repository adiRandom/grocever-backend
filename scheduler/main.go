package main

import (
	"context"
	"scheduler/gateways/events"
	"scheduler/services/crawl"
)

func main() {
	crawlScheduler := crawl.GetScheduler()
	defer crawlScheduler.Close()

	broker := events.GetRabbitMqBroker()
	broker.Start(context.Background())
}

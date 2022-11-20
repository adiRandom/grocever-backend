package events

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"lib/data/dto/scheduling"
	"lib/events/rabbitmq"
	amqpLib "lib/network/amqp"
	"scheduler/services/crawl"
)

var broker *rabbitmq.JsonBroker[scheduling.CrawlDto]

func processJsonMessage(msg scheduling.CrawlDto,
	_ *amqp.Channel,
	_ *amqp.Queue,
	_ context.Context,
) {
	scheduler := crawl.GetScheduler()
	scheduler.ScheduleCrawl(msg)
}

func GetRabbitMqBroker() *rabbitmq.JsonBroker[scheduling.CrawlDto] {
	if broker != nil {
		return broker
	}

	broker = rabbitmq.NewJsonBroker[scheduling.CrawlDto](
		processJsonMessage,
		amqpLib.ScheduleQueue,
		nil,
		nil,
	)
	return broker
}

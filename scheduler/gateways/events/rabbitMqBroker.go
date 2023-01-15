package events

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"lib/data/dto/scheduling"
	"lib/events/rabbitmq"
	amqpLib "lib/network/amqp"
	"scheduler/services/crawl"
)

var broker *rabbitmq.JsonBroker[scheduling.CrawlScheduleDto]

func processJsonMessage(msg scheduling.CrawlScheduleDto,
	_ *amqp.Channel,
	_ *amqp.Queue,
	_ context.Context,
) {
	scheduler := crawl.GetScheduler()
	scheduler.ScheduleCrawl(msg)
}

func GetRabbitMqBroker() *rabbitmq.JsonBroker[scheduling.CrawlScheduleDto] {
	if broker != nil {
		return broker
	}

	broker = rabbitmq.NewJsonBroker[scheduling.CrawlScheduleDto](
		processJsonMessage,
		amqpLib.ScheduleQueue,
		nil,
		nil,
	)
	return broker
}

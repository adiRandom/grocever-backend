package events

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"lib/data/dto"
	"lib/events/rabbitmq"
	amqpLib "lib/network/amqp"
	"log"
	"productProcessing/services"
	"time"
)

var rabbitMqBroker *rabbitmq.JsonBroker[dto.ProductProcessDto]
var requeueTimeout = 5 * time.Minute

func processJsonMessage(msg dto.ProductProcessDto,
	scheduleCh *amqp.Channel,
	scheduleQueue *amqp.Queue,
	_ context.Context,
) {
	service := services.NewProductService(*scheduleQueue, scheduleCh, &requeueTimeout)
	errs := service.ProcessCrawlProduct(msg)

	for _, err := range errs {
		log.Fatal(err)
	}
}

func GetRabbitMqBroker() *rabbitmq.JsonBroker[dto.ProductProcessDto] {
	if rabbitMqBroker != nil {
		return rabbitMqBroker
	}

	rabbitMqBroker = rabbitmq.NewJsonBroker[dto.ProductProcessDto](
		processJsonMessage,
		amqpLib.SearchQueue,
		&amqpLib.ScheduleQueue,
		nil,
	)
	return rabbitMqBroker
}

package events

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"lib/data/dto"
	"lib/events/rabbitmq"
	amqpLib "lib/network/amqp"
	"log"
	"productProcessing/services"
)

var rabbitMqBroker *rabbitmq.JsonBroker[dto.ProductProcessDto]

func processJsonMessage(msg dto.ProductProcessDto,
	_ *amqp.Channel,
	_ *amqp.Queue,
	_ context.Context,
) {
	service := services.ProductService{}
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
		nil,
		nil,
	)
	return rabbitMqBroker
}

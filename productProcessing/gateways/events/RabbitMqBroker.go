package events

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"lib/data/dto"
	"lib/events/rabbitmq"
	amqpLib "lib/network/amqp"
	"log"
	"productProcessing/data/database/repositories"
	"productProcessing/services"
)

var rabbitMqBroker *rabbitmq.JsonBroker[dto.ProductProcessDto]

var productService = services.NewProductService(
	repositories.GetProductRepository(),
	repositories.GetOcrProductRepository(),
	repositories.GetUserProductRepository(),
)

func processJsonMessage(msg dto.ProductProcessDto,
	_ *amqp.Channel,
	_ *amqp.Queue,
	_ context.Context,
) {
	fmt.Printf("Processing message: %+v \n", msg)
	errs := productService.ProcessCrawlProduct(msg)

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
		amqpLib.ProductProcessQueue,
		&amqpLib.ScheduleQueue,
		nil,
	)
	return rabbitMqBroker
}

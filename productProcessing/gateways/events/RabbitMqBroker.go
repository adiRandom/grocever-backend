package events

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"lib/data/dto"
	"lib/events/rabbitmq"
	amqpLib "lib/network/amqp"
	"productProcessing/data/database/repositories"
	"productProcessing/services"
	"productProcessing/services/api/nlp"
	"productProcessing/services/product"
)

var rabbitMqBroker *rabbitmq.JsonBroker[dto.ProductProcessDto]

func getOnMsg(productService *product.ProductService) func(msg dto.ProductProcessDto, outCh *amqp.Channel, outQ *amqp.Queue, ctx context.Context) {
	return func(msg dto.ProductProcessDto, outCh *amqp.Channel, outQ *amqp.Queue, ctx context.Context) {
		fmt.Printf("Processing message: %+v \n", msg)
		errs := productService.ProcessCrawlProduct(msg)

		for _, err := range errs {
			fmt.Printf("%v", err)
		}
	}
}

func GetRabbitMqBroker() *rabbitmq.JsonBroker[dto.ProductProcessDto] {
	if rabbitMqBroker != nil {
		return rabbitMqBroker
	}

	productService := product.NewProductService(
		repositories.GetProductRepository(
			repositories.GetMissLinkRepository(),
			repositories.GetOcrProductRepository(
				repositories.GetMissLinkRepository(),
				services.NewNotificationService(),
			),
			nlp.GetClient(),
		),
		repositories.GetOcrProductRepository(
			repositories.GetMissLinkRepository(),
			services.NewNotificationService(),
		),
		repositories.GetUserProductRepository(),
	)

	rabbitMqBroker = rabbitmq.NewJsonBroker[dto.ProductProcessDto](
		getOnMsg(productService),
		amqpLib.ProductProcessQueue,
		&amqpLib.ScheduleQueue,
		nil,
	)
	return rabbitMqBroker
}

package events

import (
	"context"
	"fmt"
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
	fmt.Printf("Processing message: %+v \n", msg)
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
		amqpLib.ProductProcessQueue,
		&amqpLib.ScheduleQueue,
		nil,
	)
	return rabbitMqBroker
}

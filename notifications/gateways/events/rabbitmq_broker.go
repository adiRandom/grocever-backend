package events

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"lib/data/dto/messages"
	"lib/events/rabbitmq"
	amqpLib "lib/network/amqp"
	services "notifications/services"
)

var rabbitMqBroker *rabbitmq.JsonBroker[messages.NotificationDto]

func getOnMsg(
	notificationService *services.NotificationService,
) func(msg messages.NotificationDto,
	outCh *amqp.Channel,
	outQ *amqp.Queue,
	ctx context.Context,
) {
	return func(msg messages.NotificationDto, outCh *amqp.Channel, outQ *amqp.Queue, ctx context.Context) {
		fmt.Printf("Processing message: %+v \n", msg)
		notificationService.SendNotification(msg.UserIds, ctx)
	}
}

func GetRabbitMqBroker(notificationService *services.NotificationService) *rabbitmq.JsonBroker[messages.NotificationDto] {
	if rabbitMqBroker != nil {
		return rabbitMqBroker
	}

	rabbitMqBroker = rabbitmq.NewJsonBroker[messages.NotificationDto](
		getOnMsg(notificationService),
		amqpLib.NotificationQueue,
		nil,
		nil,
	)
	return rabbitMqBroker
}

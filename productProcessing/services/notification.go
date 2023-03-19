package services

import (
	"fmt"
	"lib/events/rabbitmq"
	"lib/network/amqp"
)

type NotificationService struct {
}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func (s *NotificationService) SendNotification(userIds []uint) {
	err := rabbitmq.PushToQueue(amqp.NotificationQueue, userIds)
	if err != nil {
		fmt.Println(err)
		return
	}
}

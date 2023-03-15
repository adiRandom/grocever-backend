package main

import (
	"context"
	"lib/data/dto/messages"
	"lib/microservice"
	"notifications/data/database/entities"
	"notifications/data/database/repository"
	"notifications/gateways/events"
	"notifications/services"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func main() {
	opt := option.WithCredentialsFile("./serviceAccount.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic(err)
	}

	ms := microservice.AsyncMicroservice[messages.NotificationDto]{
		Microservice: microservice.Microservice{
			HasEnv:     true,
			ApiPortEnv: "API_PORT",
			GetRouter:  api.GetBaseRouter,
			DbEntities: []interface{}{entities.NotificationUser{}},
		},
		MessageBroker: *events.GetRabbitMqBroker(services.NewNotificationService(app, repository.GetNotificationUserRepository())),
	}
	ms.Start()
}

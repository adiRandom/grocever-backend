package services

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"log"
	"notifications/data/database/repository"
)

type NotificationService struct {
	app        *firebase.App
	repository *repository.NotificationUserRepository
}

func NewNotificationService(
	app *firebase.App,
	repository *repository.NotificationUserRepository,
) *NotificationService {
	return &NotificationService{
		app:        app,
		repository: repository,
	}
}

func (s *NotificationService) getDefaultNotificationData() map[string]string {
	return map[string]string{
		"type": "price_update",
	}
}

func (s *NotificationService) SendNotification(userIds []int, ctx context.Context) {
	fcmTokens, err := s.repository.GetTokensByUserIds(userIds)
	if err != nil {
		fmt.Println(err)
		return
	}

	client, err := s.app.Messaging(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	message := &messaging.MulticastMessage{
		Data: map[string]string{
			"score": "850",
			"time":  "2:45",
		},
		Tokens: fcmTokens,
	}

	br, err := client.SendMulticast(context.Background(), message)
	if err != nil {
		fmt.Println(err)
		return
	}

	if br.FailureCount > 0 {
		var failedTokens []string
		for idx, resp := range br.Responses {
			if !resp.Success {
				// The order of responses corresponds to the order of the registration tokens.
				failedTokens = append(failedTokens, fcmTokens[idx])
			}
		}

		fmt.Printf("List of tokens that caused failures: %v\n", failedTokens)
	} else {
		fmt.Println("All notifications sent successfully")
	}
}

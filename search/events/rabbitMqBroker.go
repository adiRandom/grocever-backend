package events

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"lib/data/dto/product"
	"lib/data/dto/scheduling"
	"lib/events/rabbitmq"
	amqpLib "lib/network/amqp"
	"search/data/repositories"
	"search/services"
	"time"
)

var rabbitMqBroker *rabbitmq.JsonBroker[product.PurchaseInstalmentDto]
var searchRequestTimeout = 1 * time.Minute

func processJsonMessage(msg product.PurchaseInstalmentDto,
	outCh *amqp.Channel,
	outQ *amqp.Queue,
	ctx context.Context,
) {
	searchService := services.NewGoogleSearchService(repositories.GetStoreMetadata())

	searchRes, err := searchService.SearchCrawlSources(msg.OcrName)
	if err != nil {
		fmt.Printf("Failed to query google for %s from store %d. Error: %s", msg.OcrName, msg.Store.StoreId, err.Error())
	}

	body := scheduling.NewCrawlScheduleDto(msg, searchRes, scheduling.Prioritized)

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("Failed to marshal search product dto for product %s from store %d. Error: %s",
			err.Error(),
			msg.OcrName,
			msg.Store.StoreId,
		)
	}

	err = outCh.PublishWithContext(ctx,
		"",        // exchange
		outQ.Name, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bodyBytes,
		})

	fmt.Printf(" [x] Sent %+v\n", body)

	if err != nil {
		fmt.Printf("Failed to publish a message to the priority crawl queue. Payload: %v. Error: %s",
			body,
			err.Error())
	}
}

func GetRabbitMqBroker() *rabbitmq.JsonBroker[product.PurchaseInstalmentDto] {
	if rabbitMqBroker != nil {
		return rabbitMqBroker
	}

	rabbitMqBroker = rabbitmq.NewJsonBroker[product.PurchaseInstalmentDto](
		processJsonMessage,
		amqpLib.SearchQueue,
		&amqpLib.ScheduleQueue,
		&searchRequestTimeout,
	)
	return rabbitMqBroker
}

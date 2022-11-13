package events

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"lib/data/dto"
	"lib/data/dto/scheduling"
	"lib/events/rabbitmq"
	amqpLib "lib/network/amqp"
	"log"
	"search/services"
	"time"
)

var rabbitMqBroker *rabbitmq.JsonBroker[dto.OcrProductDto]
var searchRequestTimeout = 1 * time.Minute

func processJsonMessage(msg dto.OcrProductDto,
	outCh *amqp.Channel,
	outQ *amqp.Queue,
	ctx context.Context,
) {
	searchService := services.GoogleSearchService{}

	searchRes, err := searchService.SearchCrawlSources(msg.ProductName)
	if err != nil {
		log.Fatalf("Failed to query google for %s from store %d. Error: %s", msg.ProductName, msg.StoreId, err.Error())
	}

	body := scheduling.CrawlDto{
		Product: dto.CrawlProductDto{
			OcrProduct:   msg,
			CrawlSources: searchRes,
		},
		Type: scheduling.Prioritized,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("Failed to marshal search product dto for product %s from store %d. Error: %s",
			err.Error(),
			msg.ProductName,
			msg.StoreId)
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

	if err != nil {
		log.Fatalf("Failed to publish a message to the priority crawl queue. Payload: %s. Error: %s",
			body.String(),
			err.Error())
	}
}

func GetRabbitMqBroker() *rabbitmq.JsonBroker[dto.OcrProductDto] {
	if rabbitMqBroker != nil {
		return rabbitMqBroker
	}

	rabbitMqBroker = rabbitmq.NewJsonBroker[dto.OcrProductDto](
		processJsonMessage,
		amqpLib.SearchQueue,
		&amqpLib.ScheduleQueue,
		&searchRequestTimeout,
	)
	return rabbitMqBroker
}

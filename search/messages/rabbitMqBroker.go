package messages

import (
	"context"
	"dealScraper/lib/data/dto"
	"dealScraper/lib/messages"
	"dealScraper/lib/network"
	"dealScraper/search/services"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

var rabbitMqBroker *messages.RabbitMqJsonBroker[dto.OcrProductDto]
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

	body := dto.SearchProductDto{
		OcrProduct:   msg,
		CrawlSources: searchRes,
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

func GetRabbitMqBroker() *messages.RabbitMqJsonBroker[dto.OcrProductDto] {
	if rabbitMqBroker != nil {
		return rabbitMqBroker
	}

	rabbitMqBroker = &messages.RabbitMqJsonBroker[dto.OcrProductDto]{}
	rabbitMqBroker.InboundQueueName = network.SearchQueue
	rabbitMqBroker.OutboundQueueName = network.PriorityCrawlQueue
	rabbitMqBroker.ProcessTimeout = &searchRequestTimeout
	rabbitMqBroker.ProcessJsonMessage = processJsonMessage

	return rabbitMqBroker
}

package messages

import (
	crawlers "dealScraper/crawlers/services"
	"dealScraper/lib/data/dto"
	"dealScraper/lib/messages"
	"dealScraper/lib/network"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

var rabbitMqBroker *messages.RabbitMqJsonMultiplexBroker[dto.SearchProductDto]
var messageProcessingTimeout = 1 * time.Minute
var inboundQueues = map[string]messages.AmqpJsonMultiplexBrokerInboundQueueMetadata{
	network.PriorityCrawlQueue: {
		QueueName:      network.PriorityCrawlQueue,
		ProcessTimeout: &messageProcessingTimeout,
	},
	network.CrawlQueue: {
		QueueName:      network.CrawlQueue,
		ProcessTimeout: &messageProcessingTimeout,
	},
}

const queueSwitchInterval = 10

func pickInboundQueue(currentQueueName string, queueMetadata messages.AmqpJsonMultiplexBrokerSelectQueueMetadataMap[dto.SearchProductDto]) string {
	if currentQueueName == "" {
		return network.CrawlQueue
	}

	if currentQueueName == network.PriorityCrawlQueue {
		if queueMetadata[network.PriorityCrawlQueue].MessageCount == 0 {
			return network.CrawlQueue
		}

		if queueMetadata[network.PriorityCrawlQueue].DeltaProcessedCount >= queueSwitchInterval {
			return network.CrawlQueue
		}
	}

	if currentQueueName == network.CrawlQueue {
		if queueMetadata[network.PriorityCrawlQueue].MessageCount == 0 {
			return network.CrawlQueue
		}

		if queueMetadata[network.CrawlQueue].MessageCount == 0 {
			println("Reason: No messages in queue")
			return network.PriorityCrawlQueue
		}

		println("Delta processed count: ", queueMetadata[network.CrawlQueue].DeltaProcessedCount)

		if queueMetadata[network.CrawlQueue].DeltaProcessedCount >= queueSwitchInterval {
			println("Reason: Switch interval reached")
			return network.PriorityCrawlQueue
		}
	}

	return currentQueueName
}

func processJsonMessage(args messages.AmqpJsonMultiplexBrokerProcessArgs[dto.SearchProductDto]) {
	crawlRes := crawlers.CrawlProductPages(args.Msg.CrawlSources)
	body := dto.ProductProcessDto{OcrProductDto: args.Msg.OcrProduct, CrawlResults: crawlRes}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("Failed to marshal process product dto for product %s from store %d. Error: %s",
			err.Error(),
			body.OcrProductDto.ProductName,
			body.OcrProductDto.ProductName)
	}

	err = args.OutCh.PublishWithContext(args.Ctx,
		"",             // exchange
		args.OutQ.Name, // routing key
		false,          // mandatory
		false,          // immediate
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

func GetRabbitMqBroker() *messages.RabbitMqJsonMultiplexBroker[dto.SearchProductDto] {
	if rabbitMqBroker != nil {
		return rabbitMqBroker
	}

	rabbitMqBroker = messages.NewRabbitMqJsonMultiplexBroker[dto.SearchProductDto](
		processJsonMessage,
		inboundQueues,
		network.ProductProcessQueue,
		pickInboundQueue,
		nil,
	)
	return rabbitMqBroker
}

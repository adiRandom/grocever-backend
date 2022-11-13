package events

import (
	crawlers "crawlers/services"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"lib/data/dto"
	"lib/events/rabbitmq/multiplex"
	amqpLib "lib/network/amqp"
	"log"
	"time"
)

var rabbitMqBroker *multiplex.JsonBroker[dto.CrawlProductDto]
var messageProcessingTimeout = 1 * time.Minute
var inboundQueues = multiplex.InQueues{
	amqpLib.PriorityCrawlQueue: *multiplex.NewInQueueMetadata(
		amqpLib.PriorityCrawlQueue,
		&messageProcessingTimeout,
	),
	amqpLib.CrawlQueue: *multiplex.NewInQueueMetadata(
		amqpLib.CrawlQueue,
		&messageProcessingTimeout,
	),
}

const queueSwitchInterval = 10

func pickInboundQueue(currentQueueName string,
	queueMetadata multiplex.OnSelectQueueCtx[dto.CrawlProductDto],
) string {
	if currentQueueName == "" {
		return amqpLib.CrawlQueue
	}

	if currentQueueName == amqpLib.PriorityCrawlQueue {
		if queueMetadata[amqpLib.PriorityCrawlQueue].MessageCount == 0 {
			return amqpLib.CrawlQueue
		}

		if queueMetadata[amqpLib.PriorityCrawlQueue].DeltaProcessedCount >= queueSwitchInterval {
			return amqpLib.CrawlQueue
		}
	}

	if currentQueueName == amqpLib.CrawlQueue {
		if queueMetadata[amqpLib.PriorityCrawlQueue].MessageCount == 0 {
			return amqpLib.CrawlQueue
		}

		if queueMetadata[amqpLib.CrawlQueue].MessageCount == 0 {
			println("Reason: No messages in queue")
			return amqpLib.PriorityCrawlQueue
		}

		println("Delta processed count: ", queueMetadata[amqpLib.CrawlQueue].DeltaProcessedCount)

		if queueMetadata[amqpLib.CrawlQueue].DeltaProcessedCount >= queueSwitchInterval {
			println("Reason: Switch interval reached")
			return amqpLib.PriorityCrawlQueue
		}
	}

	return currentQueueName
}

func processJsonMessage(args multiplex.OnMessageArgs[dto.CrawlProductDto]) {
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

func GetRabbitMqBroker() *multiplex.JsonBroker[dto.CrawlProductDto] {
	if rabbitMqBroker != nil {
		return rabbitMqBroker
	}

	rabbitMqBroker = multiplex.NewJsonBroker[dto.CrawlProductDto](
		processJsonMessage,
		inboundQueues,
		amqpLib.ProductProcessQueue,
		pickInboundQueue,
		nil,
	)
	return rabbitMqBroker
}

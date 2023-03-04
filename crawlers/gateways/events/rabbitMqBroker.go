package events

import (
	crawlers "crawlers/services"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"lib/data/dto"
	"lib/data/dto/crawl"
	crawlModels "lib/data/models/crawl"
	"lib/events/rabbitmq/multiplex"
	"lib/functional"
	amqpLib "lib/network/amqp"
	"time"
)

var rabbitMqBroker *multiplex.JsonBroker[crawl.ProductDto]
var messageProcessingTimeout = 1 * time.Minute
var deadlockTimeout = 5 * time.Minute
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
	queueMetadata multiplex.OnSelectQueueCtx[crawl.ProductDto],
) string {
	if currentQueueName == "" {
		return amqpLib.PriorityCrawlQueue
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

func processJsonMessage(args multiplex.OnMessageArgs[crawl.ProductDto]) {
	println("Processing message from queue: ", args.From)

	crawlRes := crawlers.CrawlProductPages(args.Msg.CrawlSources)
	if crawlRes == nil {
		return
	}

	body := dto.ProductProcessDto{
		OcrProduct: args.Msg.OcrProduct,
		CrawlResults: functional.Map(crawlRes, func(res crawlModels.ResultModel) crawl.ResultDto {
			return res.ToDto()
		}),
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("Failed to marshal process product dto for product %s from store %d. Error: %s", args.Msg.OcrProduct.OcrName, args.Msg.OcrProduct.Store.StoreId, err.Error())
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
		fmt.Printf("Failed to publish a message to the priority crawl queue. Payload: %v. Error: %s",
			body,
			err.Error())
	}
}

func GetRabbitMqBroker() *multiplex.JsonBroker[crawl.ProductDto] {
	if rabbitMqBroker != nil {
		return rabbitMqBroker
	}

	rabbitMqBroker = multiplex.NewJsonBroker[crawl.ProductDto](
		processJsonMessage,
		inboundQueues,
		amqpLib.ProductProcessQueue,
		pickInboundQueue,
		&deadlockTimeout,
	)
	return rabbitMqBroker
}

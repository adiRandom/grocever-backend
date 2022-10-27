package services

import (
	"context"
	"dealScraper/lib/data/dto"
	"dealScraper/lib/helpers"
	"dealScraper/lib/network"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

const searchRequestTimeout = 1 * time.Minute

func ListenAndProcessSearchRequests(ctx context.Context) {
	conn, ch, q, connErr := network.GetRabbitMQConnection(network.SearchQueue)
	var crawlConn, crawlCh, crawlQ, crawlConnErr = network.GetRabbitMQConnection(network.PriorityCrawlQueue)

	if connErr != nil {
		helpers.PanicOnError(connErr, connErr.Reason)
	}
	defer helpers.SafeClose(conn)
	defer helpers.SafeClose(ch)

	if crawlConnErr != nil {
		helpers.PanicOnError(crawlConnErr, crawlConnErr.Reason)
	}
	defer helpers.SafeClose(crawlConn)
	defer helpers.SafeClose(crawlCh)

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	helpers.PanicOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go processMessages(msgs, crawlCh, crawlQ)

	select {
	case <-ctx.Done():
		{
			return
		}
	case <-forever:
		{
			return
		}
	}
}

func processMessages(
	msgs <-chan amqp.Delivery,
	ch *amqp.Channel,
	q *amqp.Queue,
) {
	for msg := range msgs {
		log.Printf("Received a message: %s", msg.Body)

		var ocrProduct dto.OcrProductDto
		err := json.Unmarshal(msg.Body, &ocrProduct)
		if err != nil {
			log.Fatalf("Failed to unmarshal message. Error: %s", err.Error())
		}

		ctx, _ := context.WithTimeout(
			context.Background(),
			searchRequestTimeout,
		)

		processSearchRequest(ocrProduct, ch, q, ctx)
	}
}

func processSearchRequest(
	ocrProduct dto.OcrProductDto,
	ch *amqp.Channel,
	q *amqp.Queue,
	ctx context.Context,
) {
	searchService := GoogleSearchService{}

	searchRes, err := searchService.SearchCrawlSources(ocrProduct.ProductName)
	if err != nil {
		log.Fatalf("Failed to query google for %s from store %d. Error: %s", ocrProduct.ProductName, ocrProduct.StoreId, err.Error())
	}

	body := dto.SearchProductDto{
		OcrProduct:   ocrProduct,
		CrawlSources: searchRes,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("Failed to marshal search product dto for product %s from store %d. Error: %s",
			err.Error(),
			ocrProduct.ProductName,
			ocrProduct.StoreId)
	}

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
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

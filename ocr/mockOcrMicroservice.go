package ocr

import (
	"context"
	amqpLib "dealScraper/lib/amqp"
	"dealScraper/lib/data/dto"
	"dealScraper/lib/helpers"
	"encoding/json"
	"log"
	"time"
)
import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func SendOcrProductToQueue() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	helpers.PanicOnError(err, "Failed to connect to RabbitMQ")
	defer helpers.SafeClose(conn)

	ch, err := conn.Channel()
	helpers.PanicOnError(err, "Failed to open a channel")
	defer helpers.SafeClose(ch)

	q, err := ch.QueueDeclare(
		amqpLib.SearchQueue,
		false,
		false,
		false,
		false,
		nil,
	)
	helpers.PanicOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := dto.OcrProductDto{
		ProductName:  "test",
		ProductPrice: 1.0,
		StoreId:      1,
	}

	bodyBytes, err := json.Marshal(body)
	helpers.PanicOnError(err, "Failed to marshal body")

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bodyBytes,
		})
	helpers.PanicOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}

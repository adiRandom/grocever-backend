package main

import (
	"context"
	"encoding/json"
	"lib/data/dto"
	"lib/helpers"
	amqpLib "lib/network/amqp"
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
		ProductName:  "Lapte zuzu",
		ProductPrice: 10,
		StoreId:      nil,
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

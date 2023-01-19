package events

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"lib/data/dto/product"
	"lib/helpers"
	amqpLib "lib/network/amqp"
	"os"
)

func PublishMockProducts() {
	conn, ch, q, err := amqpLib.GetConnection(&amqpLib.SearchQueue)
	if err != nil {
		panic(err)
	}

	defer helpers.SafeClose(conn)
	defer helpers.SafeClose(ch)

	// Read the UserOcrProductDto array from products.json and publish it to the queue

	dat, err2 := os.ReadFile("products.json")
	if err2 != nil {
		panic(err2)
	}

	println(string(dat))

	var dtos []product.UserOcrProductDto
	err2 = json.Unmarshal(dat, &dtos)
	if err2 != nil {
		panic(err2)
	}

	for _, dto := range dtos {

		body, err := json.Marshal(dto)
		if err != nil {
			panic(err)
		}
		err2 = ch.PublishWithContext(context.Background(),
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)
		if err != nil {
			panic(err)
		}
	}
}

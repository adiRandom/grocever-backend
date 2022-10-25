package network

import (
	amqpLib "dealScraper/lib/amqp"
	"dealScraper/lib/helpers"
	amqp "github.com/rabbitmq/amqp091-go"
)

func GetRabbitMQConnection() (*amqp.Connection, *amqp.Channel, *amqp.Queue, *helpers.Error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, nil, nil, &helpers.Error{Msg: "Failed to connect to RabbitMQ", Reason: err.Error()}
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, nil, &helpers.Error{Msg: "Failed to open a channel", Reason: err.Error()}
	}

	q, err := ch.QueueDeclare(
		amqpLib.SearchQueue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, nil, nil, &helpers.Error{Msg: "Failed to declare a queue", Reason: err.Error()}
	}

	return conn, ch, &q, nil
}

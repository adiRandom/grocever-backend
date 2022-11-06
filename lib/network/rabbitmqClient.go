package network

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"lib/helpers"
)

const SearchQueue = "search"
const CrawlQueue = "crawl"
const PriorityCrawlQueue = "priorityCrawl"
const ProductProcessQueue = "productProcess"

func GetRabbitMQConnection(queueName string) (*amqp.Connection, *amqp.Channel, *amqp.Queue, *helpers.Error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, nil, nil, &helpers.Error{Msg: "Failed to connect to RabbitMQ", Reason: err.Error()}
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, nil, &helpers.Error{Msg: "Failed to open a channel", Reason: err.Error()}
	}

	q, err := ch.QueueDeclare(
		queueName,
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

func GetRabbitMQConnectionWithMultipleChannels(
	queueNames []string,
) (*amqp.Connection,
	[]*amqp.Channel,
	[]*amqp.Queue,
	func(),
	*helpers.Error,
) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil,
			nil,
			nil,
			func() {},
			&helpers.Error{Msg: "Failed to connect to RabbitMQ", Reason: err.Error()}
	}

	channels := make([]*amqp.Channel, len(queueNames))
	queues := make([]*amqp.Queue, len(queueNames))

	cleanup := func() {
		for _, channel := range channels {
			helpers.SafeClose(channel)
		}
		helpers.SafeClose(conn)
	}

	for i, queueName := range queueNames {
		ch, err := conn.Channel()
		if err != nil {
			cleanup()
			return nil,
				nil,
				nil,
				func() {},
				&helpers.Error{Msg: "Failed to open a channel", Reason: err.Error()}
		}

		q, err := ch.QueueDeclare(
			queueName,
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			cleanup()
			return nil,
				nil,
				nil,
				func() {},
				&helpers.Error{Msg: "Failed to declare a queue", Reason: err.Error()}
		}

		channels[i] = ch
		queues[i] = &q
	}

	return conn, channels, queues, cleanup, nil
}

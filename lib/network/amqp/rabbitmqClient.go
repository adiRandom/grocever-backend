package amqp

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"lib/helpers"
)

var SearchQueue = "search"
var CrawlQueue = "crawl"
var PriorityCrawlQueue = "priorityCrawl"
var OcrQueue = "ocr"

const ProductProcessQueue = "productProcess"
const NotificationQueue = "notification"

var ScheduleQueue = "schedule"

func declareDefaultQueue(ch *amqp.Channel, queueName string) (amqp.Queue, error) {
	return ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
}

func GetConnection(queueName *string) (*amqp.Connection, *amqp.Channel, *amqp.Queue, *helpers.Error) {
	if queueName == nil {
		return nil, nil, nil, nil
	}

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, nil, nil, &helpers.Error{Msg: "Failed to connect to RabbitMQ", Reason: err.Error()}
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, nil, &helpers.Error{Msg: "Failed to open a channel", Reason: err.Error()}
	}

	q, err := declareDefaultQueue(ch, *queueName)

	if err != nil {
		return nil, nil, nil, &helpers.Error{Msg: "Failed to declare a queue", Reason: err.Error()}
	}

	return conn, ch, &q, nil
}

func GetConnectionWithMultipleChannels(
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

		q, err := declareDefaultQueue(ch, queueName)
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

func GetMessageCount(queueName string, ch *amqp.Channel) (int, *helpers.Error) {
	q, err := ch.QueueInspect(queueName)
	if err != nil {
		fmt.Println(err)
		return 0, &helpers.Error{Msg: "Failed to inspect queue", Reason: err.Error()}
	}
	fmt.Println(q.Messages)

	return q.Messages, nil
}

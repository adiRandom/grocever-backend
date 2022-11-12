package rabbitmq

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"lib/helpers"
	amqpLib "lib/network/amqp"
	"log"
	"time"
)

type JsonBroker[T any] struct {
	onMsg func(msg T,
		outCh *amqp.Channel,
		outQ *amqp.Queue,
		ctx context.Context,
	)
	inQueueName    string
	outQueueName   *string
	processTimeout *time.Duration
}

func NewJsonBroker[T any](
	onMsg func(msg T,
		outCh *amqp.Channel,
		outQ *amqp.Queue,
		ctx context.Context,
	),
	inQueueName string,
	outQueueName *string,
	processTimeout *time.Duration,
) *JsonBroker[T] {
	return &JsonBroker[T]{
		onMsg:          onMsg,
		inQueueName:    inQueueName,
		outQueueName:   outQueueName,
		processTimeout: processTimeout,
	}
}

// Start Listen for incoming messages and process them
func (broker JsonBroker[T]) Start(
	ctx context.Context,
) {
	conn, ch, q, connErr := amqpLib.GetConnection(&broker.inQueueName)

	if connErr != nil {
		helpers.PanicOnError(connErr, connErr.Reason)
	}
	defer helpers.SafeClose(conn)
	defer helpers.SafeClose(ch)

	outConn, outCh, outQ, outConnErr := amqpLib.GetConnection(broker.outQueueName)

	if outConnErr != nil {
		helpers.PanicOnError(outConnErr, outConnErr.Reason)
	}
	defer helpers.SafeClose(outConn)
	defer helpers.SafeClose(outCh)

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

	go broker.listenForMessages(msgs, outCh, outQ)

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

func (broker JsonBroker[T]) listenForMessages(
	msgs <-chan amqp.Delivery,
	outCh *amqp.Channel,
	outQ *amqp.Queue,
) {
	for msg := range msgs {
		log.Printf("Received a message: %s", msg.Body)

		var msgBody T
		err := json.Unmarshal(msg.Body, &msgBody)
		if err != nil {
			log.Fatalf("Failed to unmarshal message. Error: %s", err.Error())
		}

		var ctx context.Context
		if broker.processTimeout != nil {
			ctx, _ = context.WithTimeout(
				context.Background(),
				*broker.processTimeout,
			)
		} else {
			ctx = context.Background()
		}

		broker.onMsg(msgBody, outCh, outQ, ctx)
	}
}

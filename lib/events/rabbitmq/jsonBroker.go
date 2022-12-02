package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
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
	inConn         *amqp.Connection
	inCh           *amqp.Channel
	inQ            *amqp.Queue
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
	conn, ch, q, connErr := amqpLib.GetConnection(&inQueueName)
	if connErr != nil {
		helpers.PanicOnError(connErr, connErr.Reason)
	}
	return &JsonBroker[T]{
		onMsg:          onMsg,
		inConn:         conn,
		inCh:           ch,
		inQ:            q,
		outQueueName:   outQueueName,
		processTimeout: processTimeout,
	}
}

// Start Listen for incoming messages and process them
func (broker JsonBroker[T]) Start(
	ctx context.Context,
) {
	defer helpers.SafeClose(broker.inConn)
	defer helpers.SafeClose(broker.inCh)

	outConn, outCh, outQ, outConnErr := amqpLib.GetConnection(broker.outQueueName)

	if outConnErr != nil {
		helpers.PanicOnError(outConnErr, outConnErr.Reason)
	}
	defer helpers.SafeClose(outConn)
	defer helpers.SafeClose(outCh)

	msgs, err := broker.inCh.Consume(
		broker.inQ.Name,
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
			fmt.Printf("Failed to unmarshal message. Error: %s", err.Error())
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

func (broker JsonBroker[T]) SendInput(body T) {
	ctx := context.Background()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, *broker.processTimeout)
	defer cancel()

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("Failed to marshal message. Error: %s", err.Error())
	}

	err = broker.inCh.PublishWithContext(
		ctxWithTimeout,
		"",
		broker.inQ.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bodyBytes,
		},
	)
}

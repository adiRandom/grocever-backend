package messages

import (
	"context"
	"dealScraper/lib/helpers"
	"dealScraper/lib/network"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type RabbitMqJsonBroker[T any] struct {
	ProcessJsonMessage func(msg T,
		outCh *amqp.Channel,
		outQ *amqp.Queue,
		ctx context.Context,
	)
	InboundQueueName  string
	OutboundQueueName string
	ProcessTimeout    *time.Duration
}

func (broker RabbitMqJsonBroker[T]) ListenAndHandleRequests(
	ctx context.Context,
) {
	conn, ch, q, connErr := network.GetRabbitMQConnection(broker.InboundQueueName)
	var outConn, outCh, outQ, outConnErr = network.GetRabbitMQConnection(broker.OutboundQueueName)

	if connErr != nil {
		helpers.PanicOnError(connErr, connErr.Reason)
	}
	defer helpers.SafeClose(conn)
	defer helpers.SafeClose(ch)

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

	go broker.processMessages(msgs, outCh, outQ)

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

func (broker RabbitMqJsonBroker[T]) processMessages(
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
		if broker.ProcessTimeout != nil {
			ctx, _ = context.WithTimeout(
				context.Background(),
				*broker.ProcessTimeout,
			)
		} else {
			ctx = context.Background()
		}

		broker.ProcessJsonMessage(msgBody, outCh, outQ, ctx)
	}
}

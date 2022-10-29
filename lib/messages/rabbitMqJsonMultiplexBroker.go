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

// RabbitMqJsonMultiplexBrokerProcessArgs /**
//   - @param {amqp.Delivery} Msg The message to process
//   - @param {string} From The name of the queue the message was received from
//   - @param {amqp.Channel} OutCh The channel to use for publishing
//   - @param {amqp.Queue} OutQ The queue to use for publishing
//   - @param {context.Context} Ctx The context to use
//
// */
type RabbitMqJsonMultiplexBrokerProcessArgs[T any] struct {
	Msg   T
	From  string
	OutCh *amqp.Channel
	OutQ  *amqp.Queue
	Ctx   context.Context
}

type RabbitMqJsonMultiplexBrokerInboundQueue struct {
	QueueName      string
	ProcessTimeout *time.Duration
}

type RabbitMqJsonMultiplexBrokerCurrentQueueMetadata[T any] struct {
	QueueName           string
	LastMessage         T
	ProcessedCount      int
	DeltaProcessedCount int
}

// RabbitMqJsonMultiplexBroker /**
//   - @param {RabbitMqJsonMultiplexBrokerProcessArgs[T] => void} ProcessJsonMessage The function to process the message
//   - @param {RabbitMqJsonMultiplexBrokerInboundQueue[]} InboundQueues The inbound queues to listen to
//   - @param {string} OutboundQueueName The outbound queue to publish to
//   - @param {RabbitMqJsonMultiplexBrokerCurrentQueueMetadata[T] => *string}
//     PickInboundQueue The function to pick the inbound queue to process from next.
//     For the first call, it will be called with an empty object to pick the first queue.
//     @returns null if you still want to pick from the current queue, or the name of the next queue otherwise
type RabbitMqJsonMultiplexBroker[T any] struct {
	ProcessJsonMessage func(args RabbitMqJsonMultiplexBrokerProcessArgs[T])
	InboundQueues      map[string]RabbitMqJsonMultiplexBrokerInboundQueue
	OutboundQueueName  string
	PickInboundQueue   func(currentQueue RabbitMqJsonMultiplexBrokerCurrentQueueMetadata[T]) *string
}

type rabbitMqConnection struct {
	conn    *amqp.Connection
	ch      *amqp.Channel
	q       *amqp.Queue
	connErr *helpers.Error
}

// openMultipleRabbitMQConnections /**
//   - @param {map[string]RabbitMqJsonMultiplexBrokerInboundQueue} queues The queues to open connections for.
//     It's a map between their names and data about them
//   - @returns {map[string]rabbitMqConnection, func()} The connections
//     and a cleanup function to close the connections
//
// */
func openMultipleRabbitMqConnections(
	queues map[string]RabbitMqJsonMultiplexBrokerInboundQueue,
) (map[string]rabbitMqConnection, func()) {
	connections := make(map[string]rabbitMqConnection)
	cleanup := func() {
		for _, connection := range connections {
			helpers.SafeClose(connection.conn)
			helpers.SafeClose(connection.ch)
		}
	}

	for _, queueData := range queues {
		var conn, ch, q, connErr = network.GetRabbitMQConnection(queueData.QueueName)

		if connErr != nil {
			cleanup()
			helpers.PanicOnError(connErr, connErr.Reason)
		}

		connections[queueData.QueueName] = rabbitMqConnection{
			conn:    conn,
			ch:      ch,
			q:       q,
			connErr: connErr,
		}
	}

	return connections, cleanup
}

type messagesMapType = map[string]<-chan amqp.Delivery

func getMessagesFromMultipleRabbitMqQueues(
	connections map[string]rabbitMqConnection,
) map[string]<-chan amqp.Delivery {
	messagesMap := make(messagesMapType)
	for queueName, connection := range connections {
		messages, err := connection.ch.Consume(
			connection.q.Name,
			"",
			true,
			false,
			false,
			false,
			nil,
		)

		if err != nil {
			helpers.PanicOnError(err, "Failed to register a consumer")
		}
		messagesMap[queueName] = messages
	}

	return messagesMap
}

func (broker RabbitMqJsonMultiplexBroker[T]) ListenAndHandleRequests(
	ctx context.Context,
) {
	inboundConnections, cleanupInboundConnections := openMultipleRabbitMqConnections(broker.InboundQueues)
	defer cleanupInboundConnections()

	var outConn, outCh, outQ, outConnErr = network.GetRabbitMQConnection(broker.OutboundQueueName)

	if outConnErr != nil {
		helpers.PanicOnError(outConnErr, outConnErr.Reason)
	}
	defer helpers.SafeClose(outConn)
	defer helpers.SafeClose(outCh)

	messagesMap := getMessagesFromMultipleRabbitMqQueues(inboundConnections)

	var forever chan struct{}

	go broker.processMessages(messagesMap, broker.InboundQueues, outCh, outQ)

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

// parseMessageAndProcess /**
//   - @param {amqp.Delivery} msg The message to parse and process
//   - @param {string} from The name of the queue the message was received from
//   - @param {*time.Duration} timeout The timeout to use for creating the context for the processing
//   - @param {amqp.Channel} outCh The channel to use for publishing
//   - @param {amqp.Queue} outQ The queue to use for publishing
//   - @returns {T} The parsed message
// Parse the message from the AMPQ body and process it using the provided function on the broker

func (broker RabbitMqJsonMultiplexBroker[T]) parseMessageAndProcess(
	msg amqp.Delivery,
	queueName string,
	timeout *time.Duration,
	outCh *amqp.Channel,
	outQ *amqp.Queue,
) T {
	log.Printf("Received a message: %s", msg.Body)

	var msgBody T
	err := json.Unmarshal(msg.Body, &msgBody)
	if err != nil {
		log.Fatalf("Failed to unmarshal message. Error: %s", err.Error())
	}

	var ctx context.Context
	if timeout != nil {
		ctx, _ = context.WithTimeout(
			context.Background(),
			*timeout,
		)
	} else {
		ctx = context.Background()
	}

	broker.ProcessJsonMessage(
		RabbitMqJsonMultiplexBrokerProcessArgs[T]{
			msgBody,
			queueName,
			outCh,
			outQ,
			ctx,
		})

	return msgBody
}

func (broker RabbitMqJsonMultiplexBroker[T]) processMessages(
	messagesMap messagesMapType,
	inboundQueues map[string]RabbitMqJsonMultiplexBrokerInboundQueue,
	outCh *amqp.Channel,
	outQ *amqp.Queue,
) {
	currentQueueName := broker.PickInboundQueue(RabbitMqJsonMultiplexBrokerCurrentQueueMetadata[T]{})

	if currentQueueName == nil {
		panic("PickInboundQueue returned nil for the first call")
	}

	processingCount := make(map[string]int)
	deltaProcessingCount := make(map[string]int)
	currentQueueMessages := messagesMap[*currentQueueName]
	for {
		for messages := range currentQueueMessages {
			msgBody := broker.parseMessageAndProcess(
				messages,
				*currentQueueName,
				inboundQueues[*currentQueueName].ProcessTimeout,
				outCh,
				outQ,
			)

			processingCount[*currentQueueName]++
			deltaProcessingCount[*currentQueueName]++

			nextQueueName := broker.PickInboundQueue(RabbitMqJsonMultiplexBrokerCurrentQueueMetadata[T]{
				QueueName:           *currentQueueName,
				LastMessage:         msgBody,
				ProcessedCount:      processingCount[*currentQueueName],
				DeltaProcessedCount: deltaProcessingCount[*currentQueueName],
			})

			if nextQueueName != nil {
				currentQueueName = nextQueueName
				currentQueueMessages = messagesMap[*currentQueueName]
				deltaProcessingCount[*currentQueueName] = 0
				// Break from the inner for to start processing the new current queue
				break
			}
		}
	}
}

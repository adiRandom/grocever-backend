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

// AmpqJsonMultiplexBrokerProcessArgs /**
//   - @param {amqp.Delivery} Msg The message to process
//   - @param {string} From The name of the queue the message was received from
//   - @param {amqp.Channel} OutCh The channel to use for publishing
//   - @param {amqp.Queue} OutQ The queue to use for publishing
//   - @param {context.Context} Ctx The context to use
//
// */
type AmpqJsonMultiplexBrokerProcessArgs[T any] struct {
	Msg   T
	From  string
	OutCh *amqp.Channel
	OutQ  *amqp.Queue
	Ctx   context.Context
}

type AmpqJsonMultiplexBrokerInboundQueueMetadata struct {
	QueueName      string
	ProcessTimeout *time.Duration
}

type AmpqJsonMultiplexBrokerSelectQueueMetadata[T any] struct {
	QueueName           string
	LastMessage         T
	ProcessedCount      int
	DeltaProcessedCount int
	MessageCount        int
}

type AmpqJsonMultiplexBrokerSelectQueueMetadataMap[T any] map[string]AmpqJsonMultiplexBrokerSelectQueueMetadata[T]

type AmpqJsonMultiplexBrokerInboundQueueMap = map[string]AmpqJsonMultiplexBrokerInboundQueueMetadata

// RabbitMqJsonMultiplexBroker /**
//   - @param {AmpqJsonMultiplexBrokerProcessArgs[T] => void} ProcessJsonMessage The function to process the message
//   - @param {AmpqJsonMultiplexBrokerInboundQueueMetadata[]} InboundQueues The inbound queues to listen to
//   - @param {string} OutboundQueueName The outbound queue to publish to
//   - @param {AmpqJsonMultiplexBrokerSelectQueueMetadata[T] => *string}
//     PickInboundQueue The function to pick the inbound queue to process from next.
//     For the first call, it will be called with an empty string and an empty map to pick the first queue.
//     Only the current queue has the LastMessage field set.
//     @returns null if you still want to pick from the current queue, or the name of the next queue otherwise
type RabbitMqJsonMultiplexBroker[T any] struct {
	ProcessJsonMessage func(args AmpqJsonMultiplexBrokerProcessArgs[T])
	InboundQueues      AmpqJsonMultiplexBrokerInboundQueueMap
	OutboundQueueName  string
	PickInboundQueue   func(
		currentQueueName string,
		queueMetadata AmpqJsonMultiplexBrokerSelectQueueMetadataMap[T],
	) *string
}

type rabbitMqConnection struct {
	conn    *amqp.Connection
	ch      *amqp.Channel
	q       *amqp.Queue
	connErr *helpers.Error
}

// openMultipleRabbitMQConnections /**
//   - @param {map[string]AmpqJsonMultiplexBrokerInboundQueueMetadata} queues The queues to open connections for.
//     It's a map between their names and data about them
//   - @returns {map[string]rabbitMqConnection, func()} The connections
//     and a cleanup function to close the connections
//
// */
func openMultipleRabbitMqConnections(
	queues map[string]AmpqJsonMultiplexBrokerInboundQueueMetadata,
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

type connectionMapType = map[string]rabbitMqConnection

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

	var forever chan struct{}

	go broker.processMessages(inboundConnections, broker.InboundQueues, outCh, outQ)

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
		AmpqJsonMultiplexBrokerProcessArgs[T]{
			msgBody,
			queueName,
			outCh,
			outQ,
			ctx,
		})

	return msgBody
}

func createSelectQueueMetadataMap[T any](
	connectionMap connectionMapType,
	currentQueueName string,
	currentQueueMessage T,
	processedCountMap map[string]int,
	deltaProcessedMap map[string]int,
) AmpqJsonMultiplexBrokerSelectQueueMetadataMap[T] {
	metadataMap := make(AmpqJsonMultiplexBrokerSelectQueueMetadataMap[T])
	for queueName, connection := range connectionMap {
		metadata := AmpqJsonMultiplexBrokerSelectQueueMetadata[T]{
			ProcessedCount:      processedCountMap[queueName],
			DeltaProcessedCount: deltaProcessedMap[queueName],
			QueueName:           queueName,
			MessageCount:        connection.q.Messages,
		}

		if queueName == currentQueueName {
			metadata.LastMessage = currentQueueMessage
		}

		metadataMap[queueName] = metadata
	}

	return metadataMap
}

func (broker RabbitMqJsonMultiplexBroker[T]) processMessages(
	inboundConnectionMap connectionMapType,
	inboundQueueData map[string]AmpqJsonMultiplexBrokerInboundQueueMetadata,
	outCh *amqp.Channel,
	outQ *amqp.Queue,
) {
	currentQueueName := broker.PickInboundQueue("", map[string]AmpqJsonMultiplexBrokerSelectQueueMetadata[T]{})

	if currentQueueName == nil {
		panic("PickInboundQueue returned nil for the first call")
	}

	processingCount := make(map[string]int)
	deltaProcessingCount := make(map[string]int)
	currentConnexion := inboundConnectionMap[*currentQueueName]
	for {
		messages, err := currentConnexion.ch.Consume(
			*currentQueueName,
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

		for msg := range messages {
			msgBody := broker.parseMessageAndProcess(
				msg,
				*currentQueueName,
				inboundQueueData[*currentQueueName].ProcessTimeout,
				outCh,
				outQ,
			)

			processingCount[*currentQueueName]++
			deltaProcessingCount[*currentQueueName]++

			nextQueueName := broker.PickInboundQueue(
				*currentQueueName,
				createSelectQueueMetadataMap[T](
					inboundConnectionMap,
					*currentQueueName,
					msgBody,
					processingCount,
					deltaProcessingCount,
				),
			)

			if nextQueueName != nil {
				currentQueueName = nextQueueName
				currentConnexion = inboundConnectionMap[*currentQueueName]
				deltaProcessingCount[*currentQueueName] = 0
				// Break from the inner for to start processing the new current queue
				break
			}
		}
	}
}

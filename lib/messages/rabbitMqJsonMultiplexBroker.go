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

// AmqpJsonMultiplexBrokerProcessArgs /**
//   - @param {amqp.Delivery} Msg The message to process
//   - @param {string} From The name of the queue the message was received from
//   - @param {amqp.Channel} OutCh The channel to use for publishing
//   - @param {amqp.Queue} OutQ The queue to use for publishing
//   - @param {context.Context} Ctx The context to use
//
// */
type AmqpJsonMultiplexBrokerProcessArgs[T any] struct {
	Msg   T
	From  string
	OutCh *amqp.Channel
	OutQ  *amqp.Queue
	Ctx   context.Context
}

type AmqpJsonMultiplexBrokerInboundQueueMetadata struct {
	QueueName      string
	ProcessTimeout *time.Duration
}

type AmqpJsonMultiplexBrokerSelectQueueMetadata[T any] struct {
	QueueName           string
	LastMessage         T
	ProcessedCount      int
	DeltaProcessedCount int
	MessageCount        int
}

type AmqpJsonMultiplexBrokerSelectQueueMetadataMap[T any] map[string]AmqpJsonMultiplexBrokerSelectQueueMetadata[T]

type AmqpJsonMultiplexBrokerInboundQueueMap = map[string]AmqpJsonMultiplexBrokerInboundQueueMetadata

// RabbitMqJsonMultiplexBroker /**
//   - @param {AmqpJsonMultiplexBrokerProcessArgs[T] => void} ProcessJsonMessage The function to process the message
//   - @param {AmqpJsonMultiplexBrokerInboundQueueMetadata[]} InboundQueues The inbound queues to listen to
//   - @param {string} OutboundQueueName The outbound queue to publish to
//   - @param {AmqpJsonMultiplexBrokerSelectQueueMetadata[T] => *string}
//     PickInboundQueue The function to pick the inbound queue to process from next.
//     For the first call, it will be called with an empty string and an empty map to pick the first queue.
//     Only the current queue has the LastMessage field set.
//   - @param PickQueueTimeout Besides picking the next queue after processing a message, it will also pick the next queue
//     after this timeout. This is to prevent a deadlock when the current queue is empty and the other queues are full.
//     @returns null if you still want to pick from the current queue, or the name of the next queue otherwise
type RabbitMqJsonMultiplexBroker[T any] struct {
	ProcessJsonMessage func(args AmqpJsonMultiplexBrokerProcessArgs[T])
	InboundQueues      AmqpJsonMultiplexBrokerInboundQueueMap
	OutboundQueueName  string
	PickInboundQueue   func(
		currentQueueName string,
		queueMetadata AmqpJsonMultiplexBrokerSelectQueueMetadataMap[T],
	) *string
	PickQueueTimeout     *time.Duration
	processedCount       map[string]int
	deltaProcessedCount  map[string]int
	currentQueueName     *string
	currentConnection    *rabbitMqConnection
	inboundConnectionMap map[string]*rabbitMqConnection
}

func NewRabbitMqJsonMultiplexBroker[T any](
	processJsonMessage func(args AmqpJsonMultiplexBrokerProcessArgs[T]),
	inboundQueues AmqpJsonMultiplexBrokerInboundQueueMap,
	outboundQueueName string,
	pickInboundQueue func(
		currentQueueName string,
		queueMetadata AmqpJsonMultiplexBrokerSelectQueueMetadataMap[T],
	) *string,
	PickQueueTimeout *time.Duration,
) *RabbitMqJsonMultiplexBroker[T] {
	return &RabbitMqJsonMultiplexBroker[T]{
		ProcessJsonMessage:   processJsonMessage,
		InboundQueues:        inboundQueues,
		OutboundQueueName:    outboundQueueName,
		PickInboundQueue:     pickInboundQueue,
		PickQueueTimeout:     PickQueueTimeout,
		processedCount:       make(map[string]int),
		deltaProcessedCount:  make(map[string]int),
		currentQueueName:     nil,
		currentConnection:    nil,
		inboundConnectionMap: make(map[string]*rabbitMqConnection),
	}
}

type rabbitMqConnection struct {
	conn    *amqp.Connection
	ch      *amqp.Channel
	q       *amqp.Queue
	connErr *helpers.Error
}

type connectionMapType = map[string]*rabbitMqConnection

// openMultipleRabbitMQConnections /**
//   - @param {map[string]AmqpJsonMultiplexBrokerInboundQueueMetadata} queues The queues to open connections for.
//     It's a map between their names and data about them
//   - @returns {map[string]rabbitMqConnection, func()} The connections
//     and a cleanup function to close the connections
//
// */
func openMultipleRabbitMqConnections(
	queues map[string]AmqpJsonMultiplexBrokerInboundQueueMetadata,
) (connectionMapType, func()) {
	connections := make(connectionMapType)
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

		connections[queueData.QueueName] = &rabbitMqConnection{
			conn:    conn,
			ch:      ch,
			q:       q,
			connErr: connErr,
		}
	}

	return connections, cleanup
}

func (broker *RabbitMqJsonMultiplexBroker[T]) ListenAndHandleRequests(
	ctx context.Context,
) {
	inboundConnections, cleanupInboundConnections := openMultipleRabbitMqConnections(broker.InboundQueues)
	broker.inboundConnectionMap = inboundConnections
	defer cleanupInboundConnections()

	var outConn, outCh, outQ, outConnErr = network.GetRabbitMQConnection(broker.OutboundQueueName)

	if outConnErr != nil {
		helpers.PanicOnError(outConnErr, outConnErr.Reason)
	}
	defer helpers.SafeClose(outConn)
	defer helpers.SafeClose(outCh)

	var forever chan struct{}

	go broker.listenForMessages(outCh, outQ)

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
// Parse the message from the Amqp body and process it using the provided function on the broker

func (broker *RabbitMqJsonMultiplexBroker[T]) parseMessageAndProcess(
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
		AmqpJsonMultiplexBrokerProcessArgs[T]{
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
) AmqpJsonMultiplexBrokerSelectQueueMetadataMap[T] {
	metadataMap := make(AmqpJsonMultiplexBrokerSelectQueueMetadataMap[T])
	for queueName, connection := range connectionMap {
		metadata := AmqpJsonMultiplexBrokerSelectQueueMetadata[T]{
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

func (broker *RabbitMqJsonMultiplexBroker[T]) onMessageReceived(
	msg amqp.Delivery,
	outCh *amqp.Channel,
	outQ *amqp.Queue,
) {
	msgBody := broker.parseMessageAndProcess(
		msg,
		*broker.currentQueueName,
		broker.InboundQueues[*broker.currentQueueName].ProcessTimeout,
		outCh,
		outQ,
	)

	// Get the next queue
	broker.processedCount[*broker.currentQueueName]++
	broker.deltaProcessedCount[*broker.currentQueueName]++

	nextQueueName := broker.PickInboundQueue(
		*broker.currentQueueName,
		createSelectQueueMetadataMap[T](
			broker.inboundConnectionMap,
			*broker.currentQueueName,
			msgBody,
			broker.processedCount,
			broker.deltaProcessedCount,
		),
	)

	if nextQueueName != nil {
		broker.currentQueueName = nextQueueName
		broker.currentConnection = broker.inboundConnectionMap[*broker.currentQueueName]
		broker.deltaProcessedCount[*broker.currentQueueName] = 0
	}
}

func (broker *RabbitMqJsonMultiplexBroker[T]) listenForMessages(
	outCh *amqp.Channel,
	outQ *amqp.Queue,
) {
	currentQueueName := broker.PickInboundQueue("", map[string]AmqpJsonMultiplexBrokerSelectQueueMetadata[T]{})

	if currentQueueName == nil {
		panic("PickInboundQueue returned nil for the first call")
	}

	broker.currentConnection = broker.inboundConnectionMap[*currentQueueName]
	//var recheckTicker *time.Ticker = nil
	//
	//if broker.PickQueueTimeout != nil {
	//	recheckTicker = time.NewTicker(*broker.PickQueueTimeout)
	//}

	for {
		messages, err := broker.currentConnection.ch.Consume(
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

		msg := <-messages
		broker.onMessageReceived(msg, outCh, outQ)
	}
}

package multiplex

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"lib/functional"
	"lib/helpers"
	amqpLib "lib/network/amqp"
	"log"
	"time"
)

// openMultipleRabbitMQConnections /**
//   - @param {map[string]InQueueMetadata} queues The queues to open connections for.
//     It's a map between their names and data about them
//   - @returns {map[string]rabbitMqConnection, func()} The connections
//     and a cleanup function to close the connections
//
// */
func openMultipleRabbitMqConnections(
	queueMap map[string]InQueueMetadata,
) (connectionMapType, func()) {
	connections := make(connectionMapType)
	queueNames := functional.Keys(queueMap)

	conn, channels, queues, cleanup, err := amqpLib.GetConnectionWithMultipleChannels(queueNames)
	if err != nil {
		helpers.PanicOnError(err, err.Reason)
	}

	for i := range queueNames {
		var queueName = queueNames[i]
		connections[queueName] = &rabbitMqConnection{
			conn: conn,
			ch:   channels[i],
			q:    queues[i],
		}
	}

	return connections, cleanup
}

// Start Listen for messages and pricess them
func (broker *JsonBroker[T]) Start(
	ctx context.Context,
) {
	inboundConnections, cleanupInboundConnections := openMultipleRabbitMqConnections(broker.inQueues)
	broker.inboundConnectionMap = inboundConnections
	defer cleanupInboundConnections()

	var outConn, outCh, outQ, outConnErr = amqpLib.GetConnection(&broker.outQueueName)

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

func (broker *JsonBroker[T]) parseMessageAndProcess(
	msg amqp.Delivery,
	queueName string,
	timeout *time.Duration,
	outCh *amqp.Channel,
	outQ *amqp.Queue,
) {
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

	broker.onMessage(
		OnMessageArgs[T]{
			msgBody,
			queueName,
			outCh,
			outQ,
			ctx,
		})

	broker.lastMessage = &msgBody
}

func createSelectQueueMetadataMap[T any](
	connectionMap connectionMapType,
	currentQueueName string,
	currentQueueMessage *T,
	processedCountMap map[string]int,
	deltaProcessedMap map[string]int,
) OnSelectQueueCtx[T] {
	metadataMap := make(OnSelectQueueCtx[T])
	for queueName, connection := range connectionMap {
		metadata := InQueueContext[T]{
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

func (broker *JsonBroker[T]) pickNextInboundQueue() {
	nextQueueName := broker.selectQueueHandler(
		broker.currentQueueName,
		createSelectQueueMetadataMap[T](
			broker.inboundConnectionMap,
			broker.currentQueueName,
			broker.lastMessage,
			broker.processedCount,
			broker.deltaProcessedCount,
		),
	)

	if nextQueueName == broker.currentQueueName {
		return
	}

	broker.currentQueueName = nextQueueName
	broker.currentConnection = broker.inboundConnectionMap[broker.currentQueueName]
	broker.deltaProcessedCount[broker.currentQueueName] = 0
}
func (broker *JsonBroker[T]) onMessageReceived(
	msg amqp.Delivery,
	outCh *amqp.Channel,
	outQ *amqp.Queue,
) {
	broker.parseMessageAndProcess(
		msg,
		broker.currentQueueName,
		broker.inQueues[broker.currentQueueName].timeout,
		outCh,
		outQ,
	)

	// Get the next queue
	broker.processedCount[broker.currentQueueName]++
	broker.deltaProcessedCount[broker.currentQueueName]++

	broker.pickNextInboundQueue()
}

func (broker *JsonBroker[T]) getMessagesToQueueNameMap() map[string]<-chan amqp.Delivery {
	result := make(map[string]<-chan amqp.Delivery)
	for queueName, connection := range broker.inboundConnectionMap {
		messages, err := connection.ch.Consume(
			queueName,
			"",
			false,
			false,
			false,
			false,
			nil,
		)

		if err != nil {
			helpers.PanicOnError(err, "Failed to register a consumer")
		}

		result[queueName] = messages
	}
	return result
}

func (broker *JsonBroker[T]) listenForMessages(
	outCh *amqp.Channel,
	outQ *amqp.Queue,
) {
	msgChannels := broker.getMessagesToQueueNameMap()
	broker.pickNextInboundQueue()

	broker.currentConnection = broker.inboundConnectionMap[broker.currentQueueName]
	var recheckTicker *time.Ticker = nil

	if broker.changeQueueTimeout != nil {
		recheckTicker = time.NewTicker(*broker.changeQueueTimeout)
	}

	for {
		msgCh := msgChannels[broker.currentQueueName]

		if recheckTicker != nil {
			select {
			case msg := <-msgCh:
				{
					broker.onMessageReceived(msg, outCh, outQ)
				}
			case <-recheckTicker.C:
				{
					broker.pickNextInboundQueue()
				}
			}
		} else {
			println("Waiting for messages")
			msg := <-msgCh
			broker.onMessageReceived(msg, outCh, outQ)
		}
	}
}

package multiplex

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

// OnMessageArgs /**
//   - @param {amqp.Delivery} Msg The message to process
//   - @param {string} From The name of the queue the message was received from
//   - @param {amqp.Channel} OutCh The channel to use for publishing
//   - @param {amqp.Queue} OutQ The queue to use for publishing
//   - @param {context.Context} Ctx The context to use
//
// */
type OnMessageArgs[T any] struct {
	Msg   T
	From  string
	OutCh *amqp.Channel
	OutQ  *amqp.Queue
	Ctx   context.Context
}

type InQueueMetadata struct {
	queueName string
	timeout   *time.Duration
}

func NewInQueueMetadata(
	queueName string,
	timeout *time.Duration,
) *InQueueMetadata {
	return &InQueueMetadata{
		queueName: queueName,
		timeout:   timeout,
	}
}

type InQueueContext[T any] struct {
	QueueName           string
	LastMessage         *T
	ProcessedCount      int
	DeltaProcessedCount int
	MessageCount        int
}

type OnSelectQueueCtx[T any] map[string]InQueueContext[T]

// InQueues a map between the queue name and its metadata
type InQueues = map[string]InQueueMetadata

// JsonBroker /**
//   - @param {OnMessageArgs[T] => void} onMessage The function to process the message
//   - @param {InQueueMetadata[]} inQueues The inbound queues to listen to
//   - @param {string} outQueueName The outbound queue to publish to
//   - @param {InQueueContext[T] => *string}
//     selectQueueHandler The function to pick the inbound queue to process from next.
//     For the first call, it will be called with an empty string and an empty map to pick the first queue.
//     Only the current queue has the LastMessage field set.
//   - @param changeQueueTimeout Besides picking the next queue after processing a message, it will also pick the next queue
//     after this timeout. This is to prevent a deadlock when the current queue is empty and the other queues are full.
//     @returns null if you still want to pick from the current queue, or the name of the next queue otherwise
type JsonBroker[T any] struct {
	onMessage          func(args OnMessageArgs[T])
	inQueues           InQueues
	outQueueName       string
	selectQueueHandler func(
		currentQueueName string,
		ctx OnSelectQueueCtx[T],
	) string
	changeQueueTimeout   *time.Duration
	processedCount       map[string]int
	deltaProcessedCount  map[string]int
	currentQueueName     string
	currentConnection    *rabbitMqConnection
	inboundConnectionMap map[string]*rabbitMqConnection
	lastMessage          *T
}

func NewJsonBroker[T any](
	onMsg func(args OnMessageArgs[T]),
	inQueues InQueues,
	outQueues string,
	selectQueueHandler func(
		currentQueueName string,
		ctx OnSelectQueueCtx[T],
	) string,
	changeQueueTimeout *time.Duration,
) *JsonBroker[T] {
	processCount := make(map[string]int)
	deltaProcessCount := make(map[string]int)

	for queueName := range inQueues {
		processCount[queueName] = 0
		deltaProcessCount[queueName] = 0
	}

	return &JsonBroker[T]{
		onMessage:            onMsg,
		inQueues:             inQueues,
		outQueueName:         outQueues,
		selectQueueHandler:   selectQueueHandler,
		changeQueueTimeout:   changeQueueTimeout,
		processedCount:       processCount,
		deltaProcessedCount:  deltaProcessCount,
		currentQueueName:     "",
		currentConnection:    nil,
		inboundConnectionMap: make(map[string]*rabbitMqConnection),
		lastMessage:          nil,
	}
}

type rabbitMqConnection struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	q    *amqp.Queue
}

type connectionMapType = map[string]*rabbitMqConnection

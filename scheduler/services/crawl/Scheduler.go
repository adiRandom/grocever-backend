package crawl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"lib/data/dto/scheduling"
	"lib/network/amqp"
	"scheduler/usecases/crawl"
)

type Scheduler struct {
	conn    *amqp091.Connection
	chs     map[string]*amqp091.Channel
	queues  map[string]*amqp091.Queue
	cleanup func()
}

var scheduler *Scheduler

const timeout = 1 // seconds

func newScheduler() *Scheduler {
	queueNames := []string{amqp.CrawlQueue, amqp.PriorityCrawlQueue}
	conn, chs, queues, cleanup, err := amqp.GetConnectionWithMultipleChannels(queueNames)
	if err != nil {
		panic(err)
	}

	chsMap := make(map[string]*amqp091.Channel)
	queueMap := make(map[string]*amqp091.Queue)

	for i, queueName := range queueNames {
		chsMap[queueName] = chs[i]
		queueMap[queueName] = queues[i]
	}

	return &Scheduler{
		conn:    conn,
		chs:     chsMap,
		queues:  queueMap,
		cleanup: cleanup,
	}
}

func GetScheduler() *Scheduler {
	if scheduler == nil {
		scheduler = newScheduler()
	}

	return scheduler
}

func (s *Scheduler) Close() {
	s.cleanup()
}

func (s *Scheduler) ScheduleCrawl(dto scheduling.CrawlDto) {
	if dto.Type == scheduling.Requeue {
		err := requeueService.Requeue(dto.Product)
		if err != nil {
			return
		}
	}

	queueName := crawl.GetQueueForPriority(dto.Type)
	if queueName == "" {
		fmt.Printf("Invalid crawl type: %s", dto.Type)
		return
	}

	ch := s.chs[queueName]
	q := s.queues[queueName]

	body := dto.Product
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("Failed to marshal crawl source: %s. Error: %s", body.String(), err.Error())
		return
	}

	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, timeout)

	err = ch.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp091.Publishing{
			Type: "application/json",
			Body: bodyBytes,
		},
	)
	if err != nil {
		fmt.Printf("Failed to publish crawl request: %s. Error: %s", body.String(), err.Error())
	}
}

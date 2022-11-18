package test

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"lib/data/dto"
	"lib/data/dto/scheduling"
	amqpLib "lib/network/amqp"
)

func ProduceRequeueMessages() {
	_, ch, q, err := amqpLib.GetConnection(&amqpLib.ScheduleQueue)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	for i := 0; i < 30; i++ {
		product := scheduling.CrawlDto{
			Type: scheduling.Requeue,
			Product: dto.CrawlProductDto{
				OcrProduct: dto.OcrProductDto{
					ProductName: "test" + string(i),
				},
				CrawlSources: make([]dto.CrawlSourceDto, 0),
			},
		}
		bodyBytes, err := json.Marshal(product)
		if err != nil {
			panic(err)
		}
		err = ch.PublishWithContext(ctx,
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        bodyBytes,
			},
		)
		if err != nil {
			return
		}
	}
}

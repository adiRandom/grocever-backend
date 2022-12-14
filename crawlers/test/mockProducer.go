package test

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"lib/data/dto"
	amqpLib "lib/network/amqp"
)

func ProduceCrawlMessages() {
	_, ch, q, err := amqpLib.GetConnection(&amqpLib.CrawlQueue)
	if err != nil {
		panic(err)
	}

	_, pCh, pQ, err := amqpLib.GetConnection(&amqpLib.PriorityCrawlQueue)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	for i := 0; i < 30; i++ {
		product := dto.CrawlProductDto{
			OcrProduct: dto.OcrProductDto{
				ProductName: "test" + string(i),
			},
			CrawlSources: make([]dto.CrawlSourceDto, 0),
		}
		bodyBytes, err := json.Marshal(product)
		if err != nil {
			panic(err)
		}
		ch.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        bodyBytes,
		})
	}

	for i := 0; i < 30; i++ {
		product := dto.CrawlProductDto{
			OcrProduct: dto.OcrProductDto{
				ProductName: "Prio" + string(i),
			},
			CrawlSources: make([]dto.CrawlSourceDto, 0),
		}
		bodyBytes, err := json.Marshal(product)
		if err != nil {
			panic(err)
		}
		pCh.PublishWithContext(ctx, "", pQ.Name, false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        bodyBytes,
		})
	}
}

package test

import (
	"encoding/json"
	"lib/data/dto/crawl"
	"lib/events/rabbitmq"
	"lib/network/amqp"
	"os"
)

//func ProduceCrawlMessages() {
//	_, ch, q, err := amqpLib.GetConnection(&amqpLib.CrawlQueue)
//	if err != nil {
//		panic(err)
//	}
//
//	_, pCh, pQ, err := amqpLib.GetConnection(&amqpLib.PriorityCrawlQueue)
//	if err != nil {
//		panic(err)
//	}
//
//	ctx := context.Background()
//	for i := 0; i < 30; i++ {
//		product := crawl.ProductDto{
//			OcrProduct: dto.OcrProductDto{
//				ProductName: "test" + string(i),
//			},
//			CrawlSources: make([]crawl.SourceDto, 0),
//		}
//		bodyBytes, err := json.Marshal(product)
//		if err != nil {
//			panic(err)
//		}
//		ch.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
//			ContentType: "application/json",
//			Body:        bodyBytes,
//		})
//	}
//
//	for i := 0; i < 30; i++ {
//		product := crawl.ProductDto{
//			OcrProduct: dto.OcrProductDto{
//				ProductName: "Prio" + string(i),
//			},
//			CrawlSources: make([]crawl.SourceDto, 0),
//		}
//		bodyBytes, err := json.Marshal(product)
//		if err != nil {
//			panic(err)
//		}
//		pCh.PublishWithContext(ctx, "", pQ.Name, false, false, amqp.Publishing{
//			ContentType: "application/json",
//			Body:        bodyBytes,
//		})
//	}
//}

func ProduceCrawlMessages(path string) {
	// Open a file and read the contents as JSON

	// Open the file
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Read the file
	decoder := json.NewDecoder(file)
	var product crawl.ProductDto
	decoder.Decode(&product)

	rabbitmq.PushToQueue(amqp.PriorityCrawlQueue, product)
}

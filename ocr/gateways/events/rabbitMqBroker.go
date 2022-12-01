package events

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"lib/data/dto"
	"lib/data/dto/ocr"
	"lib/events/rabbitmq"
	"lib/functional"
	amqpLib "lib/network/amqp"
	"ocr/models"
	"ocr/services"
	"time"
)

var broker *rabbitmq.JsonBroker[ocr.UploadDto]
var timeout = 1 * time.Minute

func processJsonMessage(msg ocr.UploadDto,
	outCh *amqp.Channel,
	outQ *amqp.Queue,
	ctx context.Context,
) {
	reader := bytes.NewReader(msg.Bytes)
	ocrService := services.GetOcrService()
	text, err := ocrService.ProcessImage(reader)

	if err != nil {
		fmt.Printf("Failed to process image. Error: %s", err.Error())
	}

	parseService := services.GetParseService()
	products, err := parseService.GetOcrProducts(text)

	if err != nil {
		fmt.Printf("Failed to parse products. Error: %s", err.Error())
	}

	dtoProducts := functional.Map(products, func(product models.OcrProduct) dto.OcrProductDto {
		return product.ToDto()
	})
	body, err := json.Marshal(dtoProducts)
	if err != nil {
		fmt.Printf("Failed to marshal ocr product dto. Error: %s", err.Error())
	}

	err = outCh.PublishWithContext(ctx,
		"",        // exchange
		outQ.Name, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return
	}
}

func GetRabbitMqBroker() *rabbitmq.JsonBroker[ocr.UploadDto] {
	if broker == nil {
		broker = rabbitmq.NewJsonBroker[ocr.UploadDto](
			processJsonMessage,
			amqpLib.OcrQueue,
			&amqpLib.SearchQueue,
			&timeout,
		)
	}

	return broker
}

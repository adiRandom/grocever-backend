package events

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"lib/data/dto/ocr"
	productDto "lib/data/dto/product"
	"lib/data/models/product"
	"lib/events/rabbitmq"
	"lib/functional"
	amqpLib "lib/network/amqp"
	"ocr/api/product_processing"
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
	products, err := parseService.GetOcrProducts(text, msg.UserId)
	if err != nil {
		fmt.Printf("Failed to parse product. Error: %s", err.Error())
	}

	productNames := functional.Map[product.UserOcrProductModel, string](
		products,
		func(productModel product.UserOcrProductModel) string {
			return productModel.OcrProduct.OcrProductName
		},
	)

	productProcessingApiClient := product_processing.GetClient()
	exists, err := productProcessingApiClient.OcrProductsExists(productNames)

	var newProducts []product.UserOcrProductModel
	if err != nil {
		fmt.Printf("Failed to check if product exist. Error: %s", err.Error())
		newProducts = products
	} else {
		newProducts = functional.IndexedFilter[product.UserOcrProductModel](
			products,
			func(index int, _ product.UserOcrProductModel) bool {
				return !exists[index]
			},
		)
	}

	dtoProducts := functional.Map(newProducts, func(product product.UserOcrProductModel) productDto.UserOcrProductDto {
		return product.ToDto()
	})

	fmt.Printf("Sending new products to product processing.\n")
	for _, dtoProduct := range dtoProducts {
		fmt.Printf("%v \n", dtoProduct)
	}

	for _, dtoProduct := range dtoProducts {
		body, err := json.Marshal(dtoProduct)
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

package services

import (
	"context"
	"dealScraper/lib/data/dto"
	"dealScraper/lib/helpers"
	"dealScraper/lib/network"
	"encoding/json"
	"log"
)

func ListenAndProcessSearchRequests(ctx context.Context) {
	conn, ch, q, connErr := network.GetRabbitMQConnection()

	if connErr != nil {
		helpers.PanicOnError(connErr, connErr.Reason)
	}
	defer helpers.SafeClose(conn)
	defer helpers.SafeClose(ch)

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	helpers.PanicOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for msg := range msgs {
			log.Printf("Received a message: %s", msg.Body)
			var ocrProduct dto.OcrProductDto
			err := json.Unmarshal(msg.Body, &ocrProduct)
			log.Fatalf("Failed to unmarshal message. Error: %s", err.Error())
		}
	}()

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

func processSearchRequest(ocrProduct dto.OcrProductDto) {
	searchRes, err := queryGoogle(ocrProduct.ProductName)
	if err != nil {
		log.Fatalf("Failed to query google for %s from store %d. Error: %s", ocrProduct.ProductName, ocrProduct.StoreId, err.Error())
	}

}

package services

import (
	vision "cloud.google.com/go/vision/apiv1"
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
)

type OcrService struct {
	client *vision.ImageAnnotatorClient
	close  func()
}

var ocrService *OcrService

func GetOcrService() *OcrService {
	if ocrService != nil {
		return ocrService
	}

	ocrService = &OcrService{}
	ctx := context.Background()

	// Creates a client.
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	ocrService.client = client
	ocrService.close = func() {
		client.Close()
	}
	return ocrService
}

func (s *OcrService) ProcessImage(file io.Reader) (*string, error) {
	ctx := context.Background()
	image, err := vision.NewImageFromReader(file)
	if err != nil {
		log.Fatalf("Failed to create image: %v", err)
	}

	annotations, err := s.client.DetectDocumentText(ctx, image, nil)
	if err != nil {
		return nil, err
	}

	// Serialize the response as JSON
	json, err := json.MarshalIndent(annotations, "", "  ")
	if err != nil {
		return nil, err
	}

	// Convert the JSON to a string
	// Create a file to write the JSON to
	resFile, err := os.OpenFile("response.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer resFile.Close()
	// Write the bytes to the file
	_, err = resFile.Write(json)

	return &annotations.Text, nil
}

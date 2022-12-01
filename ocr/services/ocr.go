package services

import (
	vision "cloud.google.com/go/vision/apiv1"
	"context"
	"fmt"
	"io"
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
		fmt.Printf("Failed to create client: %v", err)
	}
	ocrService.client = client
	ocrService.close = func() {
		client.Close()
	}
	return ocrService
}

func (s *OcrService) ProcessImage(file io.Reader) (string, error) {
	ctx := context.Background()
	image, err := vision.NewImageFromReader(file)
	if err != nil {
		fmt.Printf("Failed to create image: %v", err)
	}

	annotations, err := s.client.DetectDocumentText(ctx, image, nil)
	if err != nil {
		return "", err
	}

	return annotations.Text, nil
}

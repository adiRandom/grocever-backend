package services

import (
	"github.com/rabbitmq/amqp091-go"
	"lib/data/dto"
	entityUtils "productProcessing/data/database/entities"
	"productProcessing/data/database/repositories"
	"time"
)

type ProductService struct {
	crawlQueue     amqp091.Queue
	crawlCh        *amqp091.Channel
	requeueTimeout *time.Time
}

func NewProductService(
	crawlQueue amqp091.Queue,
	crawlCh *amqp091.Channel,
	requeueTimeout *time.Time,
) *ProductService {
	return &ProductService{
		crawlQueue:     crawlQueue,
		crawlCh:        crawlCh,
		requeueTimeout: requeueTimeout,
	}
}

func (s *ProductService) ProcessCrawlProduct(product dto.ProductProcessDto) []error {
	repo := repositories.GetProductRepository()
	entities := entityUtils.NewProductEntities(product)

	errors := make([]error, 0)

	for _, entity := range entities {
		err := repo.Create(&entity)
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
	// schedule notifications for best price
}

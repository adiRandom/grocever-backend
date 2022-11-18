package services

import (
	"context"
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"lib/data/dto"
	"lib/data/dto/scheduling"
	"lib/functional"
	entityUtils "productProcessing/data/database/entities"
	"productProcessing/data/database/repositories"
	"time"
)

type ProductService struct {
	scheduleQueue  amqp091.Queue
	scheduleCh     *amqp091.Channel
	requeueTimeout *time.Duration
}

func NewProductService(
	scheduleQueue amqp091.Queue,
	scheduleCh *amqp091.Channel,
	requeueTimeout *time.Duration,
) *ProductService {
	return &ProductService{
		scheduleQueue:  scheduleQueue,
		scheduleCh:     scheduleCh,
		requeueTimeout: requeueTimeout,
	}
}

func (s *ProductService) ProcessCrawlProduct(product dto.ProductProcessDto) []error {
	repo := repositories.GetProductRepository()
	entities := entityUtils.NewProductEntities(product)
	requeueDto := scheduling.CrawlDto{
		Type: scheduling.Requeue,
		Product: dto.CrawlProductDto{
			OcrProduct: product.OcrProductDto,
			CrawlSources: functional.Map(entities, func(entity entityUtils.ProductEntity) dto.CrawlSourceDto {
				return entity.CrawlLink.ToDto()
			}),
		},
	}

	errors := make([]error, 0)

	for _, entity := range entities {
		err := repo.Create(&entity)
		if err != nil {
			errors = append(errors, err)
		}
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, *s.requeueTimeout)
	defer cancel()

	requestDtoBody, err := json.Marshal(requeueDto)
	if err != nil {
		errors = append(errors, err)
	}

	err = s.scheduleCh.PublishWithContext(
		ctx,
		"",
		s.scheduleQueue.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        requestDtoBody,
		},
	)
	if err != nil {
		return append(errors, err)
	}

	return errors
	// schedule notifications for best price
}

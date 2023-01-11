package services

import (
	"context"
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"lib/data/dto"
	"lib/data/dto/scheduling"
	productModel "lib/data/models/product"
	"lib/functional"
	"productProcessing/data/database/repositories"
	"time"
)

type ProductService struct {
	scheduleQueue   amqp091.Queue
	scheduleCh      *amqp091.Channel
	requeueTimeout  *time.Duration
	productRepo     *repositories.ProductRepository
	ocrProductRepo  *repositories.OcrProductRepository
	userProductRepo *repositories.UserProductRepository
}

func NewProductService(
	scheduleQueue amqp091.Queue,
	scheduleCh *amqp091.Channel,
	requeueTimeout *time.Duration,
	productRepo *repositories.ProductRepository,
	ocrProductRepo *repositories.OcrProductRepository,
	userProductRepo *repositories.UserProductRepository,

) *ProductService {
	return &ProductService{
		scheduleQueue:   scheduleQueue,
		scheduleCh:      scheduleCh,
		requeueTimeout:  requeueTimeout,
		productRepo:     productRepo,
		ocrProductRepo:  ocrProductRepo,
		userProductRepo: userProductRepo,
	}
}

func (s *ProductService) ProcessCrawlProduct(productDto dto.ProductProcessDto) []error {
	products := productModel.NewProductModelsFromProcessDto(productDto)

	requeueDto := scheduling.CrawlDto{
		Type: scheduling.Requeue,
		Product: dto.CrawlProductDto{
			OcrProduct: productDto.OcrProductDto,
			CrawlSources: functional.Map(products, func(product *productModel.Model) dto.CrawlSourceDto {
				return product.CrawlLink.ToCrawlSourceDto()
			}),
		},
	}

	errors := make([]error, 0)

	// First create all the OCR products in the DB
	var ocrProduct *productModel.OcrProductModel = nil

	if len(products) > 0 {
		ocrProduct = products[0].OcrProducts[0]
	}

	if ocrProduct != nil {
		err := s.ocrProductRepo.Create(*ocrProduct)
		if err != nil {
			errors = append(errors, err)
		}

		if len(errors) > 0 {
			return errors
		}

		// Then create all the products in the DB
		for _, product := range products {
			err := s.productRepo.Create(product, ocrProduct.OcrProductName)
			if err != nil {
				errors = append(errors, err)
			}
		}

		if len(errors) > 0 {
			return errors
		}

		userOcrProduct := productModel.NewUserOcrProductModel(
			-1,                                 // ID
			productDto.OcrProductDto.Qty,       // Qty
			productDto.OcrProductDto.Price,     // Price
			productDto.UserId,                  // UserId
			*ocrProduct,                        // OcrProduct
			productDto.OcrProductDto.UnitPrice, // UnitPrice
			uint(productDto.OcrProductDto.Store.StoreId), // StoreId
			productDto.OcrProductDto.UnitType,            // UnitType
		)
		err = s.userProductRepo.CreateModel(*userOcrProduct)
		if err != nil {
			return []error{err}
		}

		errs := s.ocrProductRepo.UpdateBestPrice(ocrProduct.OcrProductName)
		if len(errs) > 0 {
			return errs
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

	if len(errors) > 0 {
		return errors
	}

	return nil
	// schedule notifications for best price
}

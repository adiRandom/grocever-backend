package services

import (
	"lib/data/dto"
	"lib/data/models"
	productModel "lib/data/models/product"
	"productProcessing/data/database/repositories"
	"time"
)

type ProductService struct {
	requeueTimeout  *time.Duration
	productRepo     *repositories.ProductRepository
	ocrProductRepo  *repositories.OcrProductRepository
	userProductRepo *repositories.PurchaseInstalmentRepository
}

func NewProductService(
	productRepo *repositories.ProductRepository,
	ocrProductRepo *repositories.OcrProductRepository,
	userProductRepo *repositories.PurchaseInstalmentRepository,
) *ProductService {
	return &ProductService{
		productRepo:     productRepo,
		ocrProductRepo:  ocrProductRepo,
		userProductRepo: userProductRepo,
	}
}

func (s *ProductService) ProcessCrawlProduct(productDto dto.ProductProcessDto) []error {
	products := productModel.NewProductModelsFromProcessDto(productDto)
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

		// Now that everything is in the db we can set the best product for this ocr product
		_, updateErrList := s.ocrProductRepo.UpdateBestProduct(ocrProduct.OcrProductName)
		if len(updateErrList) > 0 {
			return updateErrList
		}

		if productDto.OcrProduct.UserId != -1 {
			userOcrProduct := productModel.NewPurchaseInstalmentModel(
				-1,                              // ID
				productDto.OcrProduct.Qty,       // Qty
				productDto.OcrProduct.Price,     // Price
				productDto.OcrProduct.UserId,    // UserId
				*ocrProduct,                     // OcrProduct
				productDto.OcrProduct.UnitPrice, // UnitPrice
				models.NewStoreMetadataFromDto(productDto.OcrProduct.Store), // Store
				productDto.OcrProduct.UnitName,                              // UnitType
			)
			err = s.userProductRepo.CreateModel(*userOcrProduct)
			if err != nil {
				return []error{err}
			}
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
	// schedule notifications for best price
}

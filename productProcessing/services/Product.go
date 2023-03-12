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
	existingOcrProduct, err := s.ocrProductRepo.GetById(productDto.OcrProduct.OcrName)
	if err != nil {
		return []error{err}
	}

	var existingOcrProductModel *productModel.OcrProductModel
	if existingOcrProduct != nil {
		model := existingOcrProduct.ToModel(false, false)
		existingOcrProductModel = &model
	} else {
		existingOcrProductModel = nil
	}

	products := productModel.NewProductModelsFromProcessDto(productDto, existingOcrProductModel)
	errors := make([]error, 0)

	var ocrProductToUse *productModel.OcrProductModel

	if existingOcrProduct == nil && len(products) > 0 {
		// New ocr product
		newOcrProduct := products[0].OcrProducts[0]
		if newOcrProduct != nil {
			err = s.ocrProductRepo.Create(*newOcrProduct)
			if err != nil {
				return []error{err}
			}

			ocrProductToUse = newOcrProduct
		} else {
			return nil
		}
	} else {
		ocrProductToUse = existingOcrProductModel
	}

	if productDto.OcrProduct.UserId != -1 {
		userOcrProduct := productModel.NewPurchaseInstalmentModel(
			-1,                              // ID
			productDto.OcrProduct.Qty,       // Qty
			productDto.OcrProduct.Price,     // Price
			productDto.OcrProduct.UserId,    // UserId
			*ocrProductToUse,                // OcrProduct
			productDto.OcrProduct.UnitPrice, // UnitPrice
			models.NewStoreMetadataFromDto(productDto.OcrProduct.Store), // Store
			productDto.OcrProduct.UnitName,                              // UnitType
		)
		err = s.userProductRepo.CreateModel(*userOcrProduct)
		if err != nil {
			return []error{err}
		}
	}

	// Create the new products
	for _, product := range products {
		// TODO: Batch insert
		err := s.productRepo.Create(product, ocrProductToUse.OcrProductName)
		if err != nil {
			errors = append(errors, err)
		}

		if len(errors) > 0 {
			return errors
		}
	}

	// Now that everything is in the db we can set the best product for this ocr product
	updateErr := s.ocrProductRepo.UpdateBestProductAsync(ocrProductToUse.OcrProductName)
	if updateErr != nil {
		return []error{updateErr}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
	// TODO: chedule notifications for best price
}

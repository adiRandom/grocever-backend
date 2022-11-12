package services

import (
	"lib/data/dto"
	entityUtils "productProcessing/data/database/entities"
	"productProcessing/data/database/repositories"
)

type ProductService struct {
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

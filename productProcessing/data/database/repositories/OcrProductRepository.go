package repositories

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"lib/data/database"
	"lib/data/database/repositories"
	"lib/data/models/product"
	"lib/functional"
	"lib/helpers"
	"lib/types/impl"
	"productProcessing/data/database/entities"
	"sort"
)

type OcrProductRepository struct {
	repositories.DbRepository[entities.OcrProductEntity]
	missLinkRepository *MissLinkRepository
}

var ocrRepo *OcrProductRepository = nil

func GetOcrProductRepository(missLinkRepository *MissLinkRepository) *OcrProductRepository {
	if ocrRepo == nil {
		ocrRepo = &OcrProductRepository{
			missLinkRepository: missLinkRepository,
		}
		db, err := database.GetDb()
		if err != nil {
			panic(err)
		}
		ocrRepo.Db = db
	}
	return ocrRepo
}

func (r *OcrProductRepository) GetAll() ([]entities.OcrProductEntity, error) {
	var ocrProducts []entities.OcrProductEntity
	err := r.Db.Find(&ocrProducts).Error
	return ocrProducts, err
}

func (r *OcrProductRepository) GetById(ocrName string) (*entities.OcrProductEntity, error) {
	var ocrProduct entities.OcrProductEntity
	err := r.Db.First(&ocrProduct, "ocr_product_name = ?", ocrName).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ocrProduct, err
}

func (r *OcrProductRepository) GetByIdWithJoins(name string) (*entities.OcrProductEntity, error) {
	var ocrProduct entities.OcrProductEntity
	err := r.Db.
		Preload("Related").
		Preload("Products").
		Find(&ocrProduct, "ocr_product_entities.ocr_product_name = ?", name).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ocrProduct, err
}

func (r *OcrProductRepository) getRelatedOcrProductNames(ocrName string) ([]string, error) {
	var ocrProduct entities.OcrProductEntity
	err := r.Db.First(&ocrProduct, "ocr_product_name = ?", ocrName).Error
	if err != nil {
		return nil, err
	}

	var relatedOcrProducts []entities.OcrProductEntity
	err = r.Db.Model(&ocrProduct).Association("Related").Find(&relatedOcrProducts)
	if err != nil {
		return nil, err
	}

	return functional.Map(relatedOcrProducts, func(ocrProduct entities.OcrProductEntity) string {
		return ocrProduct.OcrProductName
	}), nil
}

func (r *OcrProductRepository) Save(entity entities.OcrProductEntity) error {
	return r.Db.Save(&entity).Error
}

func (r *OcrProductRepository) Delete(entity entities.OcrProductEntity) error {
	return r.Db.Delete(&entity).Error
}

func (r *OcrProductRepository) Create(model product.OcrProductModel) error {
	entity := entities.NewOcrProductEntityFromModel(model)
	return r.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(&entity).Error
}

func (r *OcrProductRepository) CreateFromProductName(name string) (*entities.OcrProductEntity, error) {
	entity := entities.OcrProductEntity{
		OcrProductName: name,
	}
	err := r.Db.FirstOrCreate(&entity).Error
	return &entity, err
}

func (r *OcrProductRepository) Exists(ocrName string) (bool, error) {
	var ocrProduct entities.OcrProductEntity
	err := r.Db.First(&ocrProduct, "ocr_product_name = ?", ocrName).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

// Function that gets an array of ocr names and returns a bool array to represent which ocr names exist in the database
func (r *OcrProductRepository) ExistsMultiple(ocrNames []string) ([]bool, error) {
	var ocrProducts []entities.OcrProductEntity
	err := r.Db.Where("ocr_product_name IN (?)", ocrNames).Find(&ocrProducts).Error
	if err != nil {
		return nil, err
	}

	var exists = make([]bool, len(ocrNames))
	for _, ocrProduct := range ocrProducts {
		for i, ocrName := range ocrNames {
			if ocrProduct.OcrProductName == ocrName {
				exists[i] = true
			}
		}
	}

	return exists, nil
}

func (r *OcrProductRepository) traverseRelatedOcrGraph(ocrName string) ([]string, error) {
	allOcrProducts := impl.NewBasicSet[string]()
	q := impl.NewQueue[string]()

	q.Push(ocrName)
	allOcrProducts.Add(ocrName)

	for !q.IsEmpty() {
		currentOcrName := q.Pop()
		relatedOcrProducts, err := r.getRelatedOcrProductNames(currentOcrName)
		if err != nil {
			return nil, err
		}
		for _, relatedOcrProduct := range relatedOcrProducts {
			if !allOcrProducts.Contains(relatedOcrProduct) {
				q.Push(relatedOcrProduct)
				allOcrProducts.Add(relatedOcrProduct)
			}
		}
	}

	return allOcrProducts.ToSlice(), nil
}

func (r *OcrProductRepository) getAllRelatedProducts(ocrNames []string) ([]*entities.ProductEntity, error) {
	var ocrProducts []entities.OcrProductEntity
	err := r.Db.
		Where("ocr_product_name IN (?)",
			ocrNames,
		).Preload("Products").Find(&ocrProducts).Error

	if err != nil {
		return nil, err
	}

	var productSet = impl.NewIdSet[uint, *entities.ProductEntity](
		func(product *entities.ProductEntity) uint {
			return product.ID
		},
	)
	for _, ocrProduct := range ocrProducts {
		productSet.AddAll(ocrProduct.Products)
	}

	return productSet.ToSlice(), nil
}

// Get the products sorted by price and pick the cheapest one for this ocr name that doesn't collide with a miss link
func (r *OcrProductRepository) pickBestProductForOcrProductAsync(
	ocrName string,
	sortedProducts []*entities.ProductEntity,
	deniedLinks ocrProductsLinksDenied,
	result chan<- bestPriceResult,
) {
	for _, productEntity := range sortedProducts {
		if !deniedLinks.IsLinkDenied(ocrName, productEntity.ID) {
			result <- bestPriceResult{
				ocrName: ocrName,
				product: productEntity,
			}
			return
		}
	}

	result <- bestPriceResult{
		ocrName: ocrName,
		product: nil,
	}
}

func (r *OcrProductRepository) UpdateBestProductAsync(ocrName string) error {
	// Traverse all the related ocr products and find the one with the highest number of products

	allOcrProducts, err := r.traverseRelatedOcrGraph(ocrName)
	if err != nil {
		return err
	}

	products, err := r.getAllRelatedProducts(allOcrProducts)
	sort.Slice(products, func(i, j int) bool {
		return products[i].Price < products[j].Price
	})

	var bestProduct *entities.ProductEntity
	for _, productEntity := range products {
		if bestProduct == nil || productEntity.Price < bestProduct.Price {
			bestProduct = productEntity
		}
	}

	bestProductResults := make(chan bestPriceResult)
	deniedLinks, err := r.missLinkRepository.getDeniedLinksForOcrProducts(allOcrProducts)
	if err != nil {
		return err
	}

	for _, ocrProduct := range allOcrProducts {
		go r.pickBestProductForOcrProductAsync(ocrProduct, products, deniedLinks, bestProductResults)
	}

	bestProductByOcrName := make(map[string]*entities.ProductEntity)
	for i := 0; i < len(allOcrProducts); i++ {
		result := <-bestProductResults
		bestProductByOcrName[result.ocrName] = result.product
	}

	// Update the best product for each ocr product
	err = r.Db.Transaction(func(tx *gorm.DB) error {
		for ocrName, bestProduct := range bestProductByOcrName {
			return tx.
				Model(&entities.OcrProductEntity{}).
				Where("ocr_product_name = ?", ocrName).
				Update("best_product_id", bestProduct.ID).
				Error
		}
		return nil
	})
	return err
}

func (r *OcrProductRepository) GetOcrProductsByNames(names []string) (map[string]entities.OcrProductEntity, error) {
	result := map[string]entities.OcrProductEntity{}
	err := r.Db.Model(&entities.OcrProductEntity{}).Where("ocr_product_name IN (?)", names).Find(&result).Error
	return result, err
}

func (r *OcrProductRepository) deleteRelated(firstOcrProduct entities.OcrProductEntity, secondOcrProduct entities.OcrProductEntity) error {
	err := r.Db.Model(&firstOcrProduct).Association("Related").Delete(&secondOcrProduct)
	if err != nil {
		return err
	}

	err = r.Db.Model(&secondOcrProduct).Association("Related").Delete(&firstOcrProduct)
	if err != nil {
		return err
	}

	return nil
}

func (r *OcrProductRepository) BreakRelatedWithoutLinkingProduct(ocrProductName string) error {
	ocrProduct, err := r.GetByIdWithJoins(ocrProductName)

	if err != nil {
		return err
	}

	if ocrProduct == nil {
		return helpers.Error{Msg: "Missing ocr product"}
	}

	productIds := functional.Map(ocrProduct.Products, func(product *entities.ProductEntity) uint {
		return product.ID
	})

	for _, relatedOcrProduct := range ocrProduct.Related {
		commonProductCount := r.Db.
			Model(&relatedOcrProduct).
			Where("id IN ?", productIds).
			Association("Products").
			Count()

		if commonProductCount == 0 {
			err = r.deleteRelated(*relatedOcrProduct, *ocrProduct)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

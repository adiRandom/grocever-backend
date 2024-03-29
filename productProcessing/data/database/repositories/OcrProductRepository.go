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
	"log"
	"math"
	"productProcessing/data/database/entities"
	"productProcessing/services"
	"sort"
)

type OcrProductRepository struct {
	repositories.DbRepository[entities.OcrProductEntity]
	missLinkRepository  *MissLinkRepository
	notificationService *services.NotificationService
}

var ocrRepo *OcrProductRepository = nil

const similarityEpsilon = 0.02

func GetOcrProductRepository(
	missLinkRepository *MissLinkRepository,
	service *services.NotificationService,
) *OcrProductRepository {
	if ocrRepo == nil {
		ocrRepo = &OcrProductRepository{
			missLinkRepository:  missLinkRepository,
			notificationService: service,
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

func (r *OcrProductRepository) getBestProductIdsMapForOcrProducts(ocrNames []string) (map[string]int, error) {
	var ocrProducts []entities.OcrProductEntity
	err := r.Db.
		Where("ocr_product_name IN (?)",
			ocrNames,
		).Preload("BestProduct").Find(&ocrProducts).Error

	if err != nil {
		return nil, err
	}

	var bestProductIdsMap = make(map[string]int)
	for _, ocrProduct := range ocrProducts {
		if ocrProduct.BestProduct != nil {
			bestProductIdsMap[ocrProduct.OcrProductName] = int(ocrProduct.BestProduct.ID)
		}
	}

	return bestProductIdsMap, nil
}

func (r *OcrProductRepository) notifyForUpdatedBestProduct(
	currentBestProductIds map[string]int,
	newBestProduct map[string]*entities.ProductEntity,
) {
	updatedOcrNames := make([]string, 0)
	for ocrName, newBestProduct := range newBestProduct {
		if currentBestProductIds[ocrName] != int(newBestProduct.ID) {
			updatedOcrNames = append(updatedOcrNames, ocrName)
		}
	}

	if len(updatedOcrNames) > 0 {
		userIds, err := r.getUserIdsToNotify(updatedOcrNames)
		if err != nil {
			log.Println(err)
			return
		}
		r.notificationService.SendNotification(userIds)
	}
}

func (r *OcrProductRepository) getAllSimilarities(ocrProductNames []string, productIds []uint) productSimilarities {
	var similarities []entities.ProductOcrProductSimilarityEntity
	var result = make(productSimilarities)
	err := r.Db.
		Where("ocr_product_name IN (?) AND product_id IN (?)",
			ocrProductNames,
			productIds,
		).Find(&similarities).Error

	if err != nil {
		return result
	}

	for _, similarity := range similarities {
		result[similarity.ProductId] = similarity.Similarity
	}

	return result
}

func (r *OcrProductRepository) UpdateBestProductAsync(ocrName string) error {
	// Traverse all the related ocr products and find the one with the highest number of products
	allOcrProducts, err := r.traverseRelatedOcrGraph(ocrName)
	if err != nil {
		return err
	}

	currentBestProductIds, err := r.getBestProductIdsMapForOcrProducts(allOcrProducts)
	if err != nil {
		return err
	}

	products, err := r.getAllRelatedProducts(allOcrProducts)
	if err != nil {
		return err
	}

	productIds := functional.Map(products, func(product *entities.ProductEntity) uint {
		return product.ID
	})

	similarities := r.getAllSimilarities(allOcrProducts, productIds)

	sort.Slice(products, func(i, j int) bool {
		similarityDelta := math.Abs(similarities[int(products[i].ID)] - similarities[int(products[j].ID)])
		if similarityDelta > similarityEpsilon {
			return similarities[int(products[i].ID)] > similarities[int(products[j].ID)]
		} else {
			return products[i].Price < products[j].Price
		}
	})

	newBestProductCh := make(chan bestPriceResult)
	deniedLinks, err := r.missLinkRepository.getDeniedLinksForOcrProducts(allOcrProducts)
	if err != nil {
		return err
	}

	for _, ocrProduct := range allOcrProducts {
		go r.pickBestProductForOcrProductAsync(ocrProduct, products, deniedLinks, newBestProductCh)
	}

	newBestProductByOcrName := make(map[string]*entities.ProductEntity)
	for i := 0; i < len(allOcrProducts); i++ {
		result := <-newBestProductCh
		newBestProductByOcrName[result.ocrName] = result.product
	}

	// Update the best product for each ocr product
	err = r.Db.Transaction(func(tx *gorm.DB) error {
		for ocrName, bestProduct := range newBestProductByOcrName {
			err := tx.
				Model(&entities.OcrProductEntity{}).
				Where("ocr_product_name = ?", ocrName).
				Update("best_product_id", bestProduct.ID).
				Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	r.notifyForUpdatedBestProduct(currentBestProductIds, newBestProductByOcrName)

	return err
}

// ocrNames are the products that just had their best product updated
// Find all the users that have a purchase instalment for one of these products with the purchased price being
// higher than the current best product price
func (r *OcrProductRepository) getUserIdsToNotify(ocrNames []string) ([]uint, error) {
	var purchaseInstalments []entities.PurchaseInstalment
	err := r.Db.
		Joins("JOIN ocr_product_entities ON ocr_product_entities.ocr_product_name = purchase_instalments.ocr_product_name_fk").
		Joins("JOIN product_entities AS best_product ON best_product.id = ocr_product_entities.best_product_id").
		Where("ocr_product_name_fk IN (?) AND purchase_instalments.price > best_product.price", ocrNames).
		Find(&purchaseInstalments).
		Error

	if err != nil {
		return nil, err
	}

	var userIds = impl.NewBasicSet[uint]()
	for _, purchaseInstalment := range purchaseInstalments {
		userIds.Add(purchaseInstalment.UserId)
	}

	return userIds.ToSlice(), nil
}

func (r *OcrProductRepository) GetOcrProductsByNames(names []string) (map[string]*entities.OcrProductEntity, error) {
	var orcProducts []*entities.OcrProductEntity
	err := r.Db.Model(&entities.OcrProductEntity{}).
		Where("ocr_product_name IN (?)", names).
		Find(&orcProducts).Error
	if err != nil {
		return nil, err
	}

	return functional.Reduce(orcProducts,
			func(acc map[string]*entities.OcrProductEntity,
				ocrProduct *entities.OcrProductEntity) map[string]*entities.OcrProductEntity {
				acc[ocrProduct.OcrProductName] = ocrProduct
				return acc
			}, make(map[string]*entities.OcrProductEntity)),
		nil
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

func (r *OcrProductRepository) breakRelatedOcrWithoutLinkingProduct(ocrProductName string) error {
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

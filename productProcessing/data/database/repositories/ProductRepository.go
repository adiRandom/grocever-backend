package repositories

import (
	"errors"
	"gorm.io/gorm"
	"lib/data/database"
	"lib/data/database/repositories"
	"lib/data/models/product"
	"log"
	"productProcessing/data/database/entities"
	"productProcessing/services/api/nlp"
	"regexp"
	"strings"
)

var productNameUnitRegex = regexp.MustCompile(` *\d*\.?\d((kg)|(l)|(g)|(buc))`)

type ProductRepository struct {
	repositories.DbRepository[entities.ProductEntity]
	missLinkRepository   *MissLinkRepository
	ocrProductRepository *OcrProductRepository
	nlpClient            *nlp.Client
}

var pr *ProductRepository = nil

func GetProductRepository(
	missLinkRepository *MissLinkRepository,
	ocrProductRepository *OcrProductRepository,
	nlpClient *nlp.Client,
) *ProductRepository {
	if pr == nil {
		pr = &ProductRepository{
			missLinkRepository:   missLinkRepository,
			ocrProductRepository: ocrProductRepository,
			nlpClient:            nlpClient,
		}
		db, err := database.GetDb()
		if err != nil {
			panic(err)
		}
		pr.Db = db
	}
	return pr
}

func (r *ProductRepository) GetAllWithCrawlLink() ([]entities.ProductEntity, error) {
	var products []entities.ProductEntity
	err := r.Db.Preload("CrawlLink").Find(&products).Error
	return products, err
}

func (r *ProductRepository) GetProductByNameAndStoreId(
	name string,
	storeId int,
	joinOcrProduct bool,
) (*entities.ProductEntity, error) {
	var product entities.ProductEntity
	var query = r.Db.Where("name = ? AND store_id = ?", name, storeId)

	if joinOcrProduct {
		query = query.Preload("OcrProducts")
	}

	err := query.First(&product).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &product, err
}

func (r *ProductRepository) updateCrawLinkUrl(product *entities.ProductEntity, url string) error {
	return r.Db.
		Model(&entities.CrawlLinkEntity{}).
		Where("product_id = ?", product.ID).
		Updates(entities.CrawlLinkEntity{
			Url: url,
		}).
		Error
}

func (r *ProductRepository) hasOcrProduct(product *entities.ProductEntity, ocrName string) (bool, error) {
	var ocrProduct *entities.OcrProductEntity

	err := r.Db.
		Model(product).
		Where(&entities.OcrProductEntity{OcrProductName: ocrName}).
		Association("OcrProducts").
		Find(&ocrProduct)
	if err != nil {
		return false, err
	}

	return ocrProduct != nil && ocrProduct.OcrProductName != "", nil
}

func (r *ProductRepository) updateProductPrice(
	product *entities.ProductEntity,
	price float32,
) error {
	return r.Db.Model(product).Updates(entities.ProductEntity{
		Price: price,
	}).Error
}

func (r *ProductRepository) createProductAndCrawlLink(
	product *product.Model,
) (*entities.ProductEntity, error) {

	entity := entities.NewProductEntityFromModel(*product)
	hasCrawlLink := entity.CrawlLink != nil

	err := r.Db.Create(entity).Error
	if err != nil {
		return nil, err
	}

	if !hasCrawlLink {
		// Create the crawl link entity now that we have the id of the product enitity
		updatedCrawlLinkModel := product.CrawlLink
		updatedCrawlLinkModel.ProductId = int(entity.ID)
		crawlLinkEntity := entities.NewCrawlLinkEntityFromModel(updatedCrawlLinkModel)
		err = r.Db.Create(crawlLinkEntity).Error
		if err != nil {
			return nil, err
		}

		entity.CrawlLink = crawlLinkEntity
	}

	return entity, nil
}

func (r *ProductRepository) processProductNameForSimilarity(productName string) string {
	lower := strings.ToLower(productName)
	// Remove the unit from the product name
	return productNameUnitRegex.ReplaceAllString(lower, "")
}

func (r *ProductRepository) linkProductAndOcrProduct(
	ocrProductName string,
	product entities.ProductEntity,
) error {
	var ocrProduct entities.OcrProductEntity

	err := r.Db.First(&ocrProduct, "ocr_product_name = ?", ocrProductName).Error
	if err != nil {
		return err
	}

	var existingOcrProducts []entities.OcrProductEntity
	err = r.Db.Model(&product).Association("OcrProducts").Find(&existingOcrProducts)

	err = r.Db.Model(&product).Association("OcrProducts").Append(&ocrProduct)
	if err != nil {
		return err
	}

	err = r.Db.Model(&ocrProduct).Association("Products").Append(&product)
	if err != nil {
		return err
	}

	// Link this ocr product to the existing ocr product
	// Then link the existing ocr product to this ocr product
	for _, existingOcrProduct := range existingOcrProducts {
		if existingOcrProduct.OcrProductName != ocrProductName {
			err = r.Db.Model(&ocrProduct).Association("Related").Append(&existingOcrProduct)
			if err != nil {
				return err
			}
			err = r.Db.Model(&existingOcrProduct).Association("Related").Append(&ocrProduct)
			if err != nil {
				return err
			}
		}
	}

	similarity := r.nlpClient.GetSimilarity(
		strings.ToLower(ocrProductName),
		r.processProductNameForSimilarity(product.Name),
	)
	similarityEntity := entities.NewProductOcrProductSimilarityEntity(ocrProduct.OcrProductName, int(product.ID), similarity)
	err = r.Db.Create(similarityEntity).Error

	return nil
}

func (r *ProductRepository) BreakProductLinkAsync(productId int, ocrProductName string) error {
	var ocrProduct entities.OcrProductEntity

	err := r.Db.First(&ocrProduct, "ocr_product_name = ?", ocrProductName).Error
	if err != nil {
		return err
	}

	var product entities.ProductEntity
	err = r.Db.First(&product, productId).Error
	if err != nil {
		return err
	}

	err = r.Db.Model(&product).Association("OcrProducts").Delete(&ocrProduct)
	if err != nil {
		return err
	}

	err = r.Db.Model(&ocrProduct).Association("Products").Delete(&product)
	if err != nil {
		return err
	}

	err = r.ocrProductRepository.breakRelatedOcrWithoutLinkingProduct(ocrProductName)
	if err != nil {
		return err
	}

	err = r.Db.Delete(&entities.ProductOcrProductSimilarityEntity{}, "product_id = ? AND ocr_product_name = ?",
		productId,
		ocrProductName,
	).Error

	if err != nil {
		return err
	}

	go func() {
		err := r.ocrProductRepository.UpdateBestProductAsync(ocrProductName)
		if err != nil {
			log.Println(err)
		}
	}()

	return nil
}

// Create Either create a new product or add a new ocr product to an existing product
func (r *ProductRepository) Create(
	product *product.Model,
	associatedOcrProductName string,
) error {
	if product.Name == "" && product.Price == 0 {
		return nil
	}

	existingProduct, err := r.GetProductByNameAndStoreId(product.Name, product.StoreId, false)
	if err != nil {
		return err
	}

	if existingProduct == nil {
		entity, err := r.createProductAndCrawlLink(product)
		if err != nil {
			return err
		}
		err = r.linkProductAndOcrProduct(associatedOcrProductName, *entity)
		if err != nil {
			return err
		}
	} else {
		isLinkDenied, err := r.missLinkRepository.IsLinkingDenied(int(existingProduct.ID), associatedOcrProductName)
		if err != nil {
			return err
		}

		if !isLinkDenied {
			err := r.linkProductAndOcrProduct(associatedOcrProductName, *existingProduct)
			if err != nil {
				return err
			}

			err = r.updateCrawLinkUrl(existingProduct, product.CrawlLink.Url)
			if err != nil {
				return err
			}

			err = r.updateProductPrice(existingProduct, product.Price)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

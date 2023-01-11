package repositories

import (
	"errors"
	"gorm.io/gorm"
	"lib/data/database"
	"lib/data/database/repositories"
	"lib/data/models/product"
	"productProcessing/data/database/entities"
)

type ProductRepository struct {
	repositories.Repository[entities.ProductEntity]
}

var pr *ProductRepository = nil

func GetProductRepository() *ProductRepository {
	if pr == nil {
		pr = &ProductRepository{}
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

// Create Either create a new product or add a new ocr product to an existing product
func (r *ProductRepository) Create(
	product *product.Model,
	associatedOcrProductName string,
) error {
	entity := entities.NewProductEntityFromModel(*product)
	ocrProductRepo := GetOcrProductRepository()
	existingProduct, err := r.GetProductByNameAndStoreId(product.Name, product.StoreId, false)
	if err != nil {
		return err
	}

	if existingProduct == nil {
		err = r.Db.Create(entity).Error
		if err != nil {
			return err
		}

		err := ocrProductRepo.LinkProductAndOcrProduct(associatedOcrProductName, entity)
		if err != nil {
			return err
		}
	} else {
		err := ocrProductRepo.LinkProductAndOcrProduct(associatedOcrProductName, *existingProduct)
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

	return nil
}

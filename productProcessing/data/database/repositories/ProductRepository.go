package repositories

import (
	"errors"
	"gorm.io/gorm"
	"productProcessing/data/database"
	"productProcessing/data/database/entities"
)

type ProductRepository struct {
	Repository[entities.ProductEntity]
}

var pr *ProductRepository = nil

func GetProductRepository() *ProductRepository {
	if pr == nil {
		pr = &ProductRepository{}
		db, err := database.GetDb()
		if err != nil {
			panic(err)
		}
		pr.db = db
	}
	return pr
}

func (r *ProductRepository) GetAll() ([]entities.ProductEntity, error) {
	var products []entities.ProductEntity
	err := r.db.Find(&products).Error
	return products, err
}

func (r *ProductRepository) GetAllWithCrawlLink() ([]entities.ProductEntity, error) {
	var products []entities.ProductEntity
	err := r.db.Preload("CrawlLink").Find(&products).Error
	return products, err
}

func (r *ProductRepository) GetById(id uint) (*entities.ProductEntity, error) {
	var product entities.ProductEntity
	err := r.db.First(&product, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &product, err
}

func (r *ProductRepository) Save(entity entities.ProductEntity) error {
	return r.db.Save(&entity).Error
}

func (r *ProductRepository) Delete(entity entities.ProductEntity) error {
	return r.db.Delete(&entity).Error
}

func (r *ProductRepository) GetProductByNameAndStoreId(
	name string,
	storeId int,
	joinOcrProduct bool,
) (*entities.ProductEntity, error) {
	var product entities.ProductEntity
	var query = r.db.Where("name = ? AND store_id = ?", name, storeId)

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
	return r.db.
		Model(&entities.CrawlLinkEntity{}).
		Where("product_id = ?", product.ID).
		Updates(entities.CrawlLinkEntity{
			Url: url,
		}).
		Error
}

func (r *ProductRepository) hasOcrProduct(product *entities.ProductEntity, ocrName string) (bool, error) {
	var ocrProduct *entities.OcrProductEntity

	err := r.db.
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
	return r.db.Model(product).Updates(entities.ProductEntity{
		Price: price,
	}).Error
}

// Create Either create a new product or add a new ocr product to an existing product
func (r *ProductRepository) Create(
	product *entities.ProductEntity,
) error {
	ocrProduct := product.OcrProducts[0]
	ocrProductRepository := GetOcrProductRepository()
	existingProduct, err := r.GetProductByNameAndStoreId(
		product.Name,
		product.StoreId,
		false,
	)
	if err != nil {
		return err
	}

	if existingProduct == nil {
		err = r.db.Create(product).Error
		if err != nil {
			return err
		}
	} else {
		hasOcrProduct, err := r.hasOcrProduct(existingProduct, ocrProduct.OcrProductName)
		if err != nil {
			return err
		}

		if !hasOcrProduct {
			err = ocrProductRepository.AddOcrProductToProduct(*ocrProduct, *existingProduct)
			if err != nil {
				return err
			}
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

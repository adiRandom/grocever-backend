package repositories

import (
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
		query = query.Preload("OcrProduct")
	}
	err := query.First(&product).Error
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

// TODO: Test
func (r *ProductRepository) hasOcrProduct(productId uint, ocrName string) (bool, error) {
	var ocrProduct *entities.OcrProductEntity

	err := r.db.
		Where("OcrProductName = ? AND id = ?", ocrName, productId).
		Association("OcrProduct").
		Find(&ocrProduct)
	if err != nil {
		return false, err
	}

	return ocrProduct != nil, nil
}

func (r *ProductRepository) updateProductPrice(
	product *entities.ProductEntity,
	price float32,
) error {
	return r.db.Model(product).Updates(entities.ProductEntity{
		Price: price,
	}).Error
}

// AddProduct Either create a new product or add a new ocr product to an existing product
func (r *ProductRepository) AddProduct(
	product *entities.ProductEntity,
	ocrProduct *entities.OcrProductEntity,
) error {
	ocrProductRepository := GetOcrProductRepository()
	existingProduct, err := r.GetProductByNameAndStoreId(
		product.Name,
		product.StoreId,
		false,
	)
	if err != nil {
		return err
	}

	if existingProduct.ID == 0 {
		err = r.Save(*product)
		if err != nil {
			return err
		}
	} else {
		hasOcrProduct, err := r.hasOcrProduct(existingProduct.ID, ocrProduct.OcrProductName)
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

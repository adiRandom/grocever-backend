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

func (r *ProductRepository) GetProductByNameAndStoreId(name string, storeId int) (*entities.ProductEntity, error) {
	var product entities.ProductEntity
	err := r.db.Where("name = ? AND store_id = ?", name, storeId).First(&product).Error
	return &product, err
}

func (r *ProductRepository) addOcrProductToProduct(product *entities.ProductEntity,
	ocrProduct *entities.OcrProductEntity,
) error {
	return r.db.Model(product).Association("OcrProduct").Append(ocrProduct)
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

// CreateOrAddAssociation Either create a new product or add a new ocr product to an existing product
func (r *ProductRepository) CreateOrUpdateExisting(
	product *entities.ProductEntity,
	ocrProduct *entities.OcrProductEntity,
) error {
	existingProduct, err := r.GetProductByNameAndStoreId(product.Name, product.StoreId)
	if err != nil {
		return err
	}

	if existingProduct.ID == 0 {
		err = r.Save(*product)
		if err != nil {
			return err
		}
	} else {
		// TODO :
		// Check if the ocr product is already associated with the product
		// If not, add it
		// Update the crawl link url
		// Update the price
	}
	return nil
}

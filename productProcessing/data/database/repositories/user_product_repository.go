package repositories

import (
	"lib/data/database"
	"lib/data/database/repositories"
	"lib/data/models/user_product"
	"lib/helpers"
	"productProcessing/api/store"
	"productProcessing/data/database/entities"
)

type UserProductRepository struct {
	// TODO: Update entity
	repositories.RepositoryWithModel[entities.ProductEntity, user_product.Model]
}

var repo *UserProductRepository = nil

func GetUserProductRepository() *UserProductRepository {
	if repo == nil {
		repo = &UserProductRepository{}
		repo.ToModel = repo.toModel
		repo.ToEntity = repo.toEntity
		db, err := database.GetDb()
		if err != nil {
			panic(err)
		}
		repo.Db = db
	}
	return repo
}

func (r *UserProductRepository) getStoreNameForId(id int) (string, error) {
	apiClient := store.GetClient()
	stores := apiClient.GetAllStores()
	for _, storeMetadata := range stores {
		if storeMetadata.StoreId == id {
			return storeMetadata.Name, nil
		}
	}
	return "", helpers.Error{Msg: "Store not found"}
}

func (r *UserProductRepository) toModel(entity entities.ProductEntity) (user_product.Model, error) {
	storeName, err := r.getStoreNameForId(int(entity.StoreId))

	if err != nil {
		return user_product.Model{}, err
	}

	return user_product.Model{
		Id:   int(entity.ID),
		Name: entity.Name,
		// TODO: Replace with ocr price
		Price:     entity.Price,
		BestPrice: entity.Price,
		Url:       entity.CrawlLink.Url,
		Store:     storeName,
	}, nil
}

func (r *UserProductRepository) toEntity(model user_product.Model) (entities.ProductEntity, error) {
	return entities.ProductEntity{}, helpers.Error{
		// TODO: Implement
		Msg: "Not implemented",
	}
}

package repositories

import (
	"lib/data/database"
	"lib/data/database/repositories"
	"lib/data/models"
	"lib/data/models/product"
	"lib/helpers"
	"productProcessing/data/database/entities"
	"productProcessing/services/api/store"
)

type UserProductRepository struct {
	repositories.RepositoryWithModel[entities.UserOcrProduct, product.UserOcrProductModel]
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

func (r *UserProductRepository) getStoreMetadataForId(id int) (models.StoreMetadata, error) {
	apiClient := store.GetClient()
	stores, err := apiClient.GetAllStores()
	if err != nil {
		return models.StoreMetadata{}, err
	}

	for _, storeMetadataDto := range stores {
		if storeMetadataDto.StoreId == id {
			return models.NewStoreMetadataFromDto(storeMetadataDto), nil
		}
	}
	return models.StoreMetadata{}, helpers.Error{Msg: "Store not found"}
}

func (r *UserProductRepository) toModel(entity entities.UserOcrProduct) (product.UserOcrProductModel, error) {
	storeMetadata, err := r.getStoreMetadataForId(int(entity.Product.StoreId))

	if err != nil {
		return product.UserOcrProductModel{}, err
	}

	return entity.ToModel(storeMetadata), nil
}

func (r *UserProductRepository) toEntity(model product.UserOcrProductModel) (entities.UserOcrProduct, error) {
	return entities.NewUserOcrProductFromModel(model), nil
}

package repositories

import (
	"lib/data/database"
	"lib/data/database/repositories"
	"lib/data/models"
	"lib/data/models/product"
	"lib/functional"
	"lib/helpers"
	"productProcessing/data/database/entities"
	productModels "productProcessing/data/models"
	"productProcessing/services/api/store"
)

type PurchaseInstalmentRepository struct {
	repositories.DbRepositoryWithModel[entities.PurchaseInstalment, product.PurchaseInstalmentModel]
}

var repo *PurchaseInstalmentRepository = nil

func GetUserProductRepository() *PurchaseInstalmentRepository {
	if repo == nil {
		repo = &PurchaseInstalmentRepository{}
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

func (r *PurchaseInstalmentRepository) getStoreMetadataForId(id int) (models.StoreMetadata, error) {
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

func (r *PurchaseInstalmentRepository) toModel(entity entities.PurchaseInstalment) (product.PurchaseInstalmentModel, error) {
	storeMetadata, err := r.getStoreMetadataForId(int(entity.StoreId))

	if err != nil {
		return product.PurchaseInstalmentModel{}, err
	}

	return entity.ToModel(storeMetadata), nil
}

func (r *PurchaseInstalmentRepository) toEntity(model product.PurchaseInstalmentModel) (*entities.PurchaseInstalment, error) {
	return entities.NewPurchaseInstalmentFromModel(model), nil
}

func (r *PurchaseInstalmentRepository) GetUserProducts(userId int) ([]productModels.UserProduct, error) {
	var purchaseInstalments []entities.PurchaseInstalment
	err := r.Db.
		Where("user_id = ?", userId).
		Preload("OcrProduct").
		Preload("BestProduct").
		Find(&purchaseInstalments).Error
	if err != nil {
		return nil, err
	}

	instalmentsGroupedByBestProduct := functional.GroupBy(
		purchaseInstalments,
		func(purchaseInstalment entities.PurchaseInstalment,
		) *entities.ProductEntity {
			return purchaseInstalment.OcrProduct.BestProduct
		})

	var userProducts []productModels.UserProduct
	for bestProduct, purchaseInstalments := range instalmentsGroupedByBestProduct {
		storeMetadata, err := r.getStoreMetadataForId(bestProduct.StoreId)
		if err != nil {
			continue
		}

		purchaseInstalmentsModels := functional.Map(
			purchaseInstalments,
			func(purchaseInstalment entities.PurchaseInstalment) product.PurchaseInstalmentModel {
				return purchaseInstalment.ToModel(storeMetadata)
			})
		userProduct := productModels.NewUserProduct(
			bestProduct.Name,
			bestProduct.Price,
			purchaseInstalmentsModels,
			uint(storeMetadata.StoreId),
			storeMetadata.Name,
			storeMetadata.Url,
		)

		userProducts = append(userProducts, *userProduct)
	}

	return userProducts, nil
}

package repositories

import (
	"lib/data/database"
	"lib/data/database/repositories"
	productDto "lib/data/dto/product"
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
	ocrProductRepository *OcrProductRepository
}

var repo *PurchaseInstalmentRepository = nil

func GetUserProductRepository() *PurchaseInstalmentRepository {
	if repo == nil {
		repo = &PurchaseInstalmentRepository{
			ocrProductRepository: GetOcrProductRepository(),
		}
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
		Preload("OcrProduct.BestProduct").
		Preload("OcrProduct.BestProduct.CrawlLink").
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

	userProducts := make([]productModels.UserProduct, 0)
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

func (r *PurchaseInstalmentRepository) CreatePurchaseInstalment(
	dto productDto.CreatePurchaseInstalmentDto,
) (*product.PurchaseInstalmentModel, error) {
	ocrProduct, err := r.ocrProductRepository.GetById(dto.OcrName)
	if err != nil {
		return nil, err
	}

	entity := entities.NewPurchaseInstalment(
		dto.UserId,
		dto.OcrName,
		*ocrProduct,
		dto.Qty,
		dto.UnitPrice,
		dto.Qty*dto.UnitPrice,
		uint(dto.Store.StoreId),
		dto.UnitName,
	)
	err = r.Create(entity)
	if err != nil {
		return nil, err
	}

	model := entity.ToModel(models.NewStoreMetadataFromDto(dto.Store))

	return &model, nil
}

func (r *PurchaseInstalmentRepository) CreatePurchaseInstalmentNoOcr(
	dto productDto.CreatePurchaseInstalmentNoOcrWithUserDto,
) (*product.PurchaseInstalmentModel, error) {
	ocrProduct, err := r.ocrProductRepository.CreateFromProductName(dto.ProductName)
	if err != nil {
		return nil, err
	}

	entity := entities.NewPurchaseInstalment(
		dto.UserId,
		dto.ProductName,
		*ocrProduct,
		dto.Qty,
		dto.UnitPrice,
		dto.Qty*dto.UnitPrice,
		dto.StoreId,
		dto.UnitName,
	)
	err = r.Create(entity)
	if err != nil {
		return nil, err
	}

	storeModel, err := r.getStoreMetadataForId(int(dto.StoreId))
	if err != nil {
		return nil, err
	}

	model := entity.ToModel(storeModel)

	return &model, nil
}

func (r *PurchaseInstalmentRepository) CreatePurchaseInstalments(
	dto productDto.CreatePurchaseInstalmentListDto,
) ([]product.PurchaseInstalmentModel, error) {
	ocrProductNames := functional.Map(
		dto.Instalments,
		func(purchaseInstalmentDto productDto.CreatePurchaseInstalmentDto) string {
			return purchaseInstalmentDto.OcrName
		},
	)

	ocrProducts, err := r.ocrProductRepository.GetOcrProductsByNames(ocrProductNames)
	if err != nil {
		return nil, err
	}

	purchaseInstalments := functional.Map(
		dto.Instalments,
		func(dto productDto.CreatePurchaseInstalmentDto) entities.PurchaseInstalment {
			return *entities.NewPurchaseInstalment(
				dto.UserId,
				dto.OcrName,
				ocrProducts[dto.OcrName],
				dto.Qty,
				dto.UnitPrice,
				dto.Qty*dto.UnitPrice,
				uint(dto.Store.StoreId),
				dto.UnitName,
			)
		},
	)

	err = r.CreateMany(purchaseInstalments)
	if err != nil {
		return nil, err
	}

	purchaseInstalmentModels := functional.IndexedMap(
		purchaseInstalments,
		func(index int, purchaseInstalment entities.PurchaseInstalment) product.PurchaseInstalmentModel {
			return purchaseInstalment.ToModel(models.NewStoreMetadataFromDto(dto.Instalments[index].Store))
		},
	)

	return purchaseInstalmentModels, nil
}

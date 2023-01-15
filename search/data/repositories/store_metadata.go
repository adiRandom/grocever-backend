package repositories

import (
	storeDto "lib/data/dto/store"
	"lib/data/models"
	"lib/functional"
	"search/services/api/store"
)

var repo *StoreMetadata

type StoreMetadata struct {
	api *store.Client
}

func GetStoreMetadata() *StoreMetadata {
	if repo == nil {
		repo = &StoreMetadata{
			api: store.GetClient(),
		}
	}

	return repo
}

func (s *StoreMetadata) GetForUrl(url string) *models.StoreMetadata {
	stores := s.api.GetAllStores()
	dto := functional.Find(stores, func(store storeDto.MetadataDto) bool {
		return store.Url == url
	})

	if dto == nil {
		return nil
	}

	storeModel := models.NewStoreMetadataFromDto(*dto)

	return &storeModel
}

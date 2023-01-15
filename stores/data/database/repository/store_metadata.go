package repository

import (
	"errors"
	"gorm.io/gorm"
	"lib/data/database"
	"lib/data/database/repositories"
	"stores/data/database/entity"
)

type StoreMetadata struct {
	repositories.DbRepository[entity.StoreMetadata]
}

var storeMetadataRepo *StoreMetadata = nil

func GetStoreMetadataRepository() *StoreMetadata {
	if storeMetadataRepo == nil {
		storeMetadataRepo = &StoreMetadata{}
		db, err := database.GetDb()
		if err != nil {
			panic(err)
		}
		storeMetadataRepo.Db = db
	}
	return storeMetadataRepo
}

func (r *StoreMetadata) GetByName(name string) (*entity.StoreMetadata, error) {
	var storeMetadata entity.StoreMetadata
	err := r.Db.First(&storeMetadata, "name = ?", name).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &storeMetadata, err
}

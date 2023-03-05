package repositories

import (
	"lib/data/database"
	"lib/data/database/repositories"
	"productProcessing/data/database/entities"
	"productProcessing/data/models"
)

const missLinkDenyLimit = 3

type MissLinkRepository struct {
	repositories.DbRepositoryWithModel[entities.MissLink, models.MissLink]
}

var missLinkRepo *MissLinkRepository = nil

func GetMissLinkRepository() *MissLinkRepository {
	if missLinkRepo == nil {
		missLinkRepo = &MissLinkRepository{}
		db, err := database.GetDb()
		if err != nil {
			panic(err)
		}
		missLinkRepo.Db = db
		missLinkRepo.ToModel = toModel
		missLinkRepo.ToEntity = toEntity
	}
	return missLinkRepo
}

func (r *MissLinkRepository) IsLinkingDenied(productId uint, ocrName string) (bool, error) {
	var missLinkCount int64 = 0
	err := r.Db.Model(&entities.MissLink{}).Where("product_id_fk = ? and ocr_product_name_fk = ?", productId, ocrName).Count(&missLinkCount).Error
	if err != nil {
		return false, err
	}

	return missLinkCount >= missLinkDenyLimit, nil
}

func toModel(entity entities.MissLink) (models.MissLink, error) {
	return *models.NewMissLinkModelFromEntity(&entity), nil
}

func toEntity(model models.MissLink) (*entities.MissLink, error) {
	return model.ToEntity(), nil
}

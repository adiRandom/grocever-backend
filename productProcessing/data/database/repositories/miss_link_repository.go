package repositories

import (
	"lib/data/database"
	"lib/data/database/repositories"
	"lib/helpers"
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

func (r *MissLinkRepository) Create(productId int, ocrName string, userId int) (*entities.MissLink, error) {
	entity := entities.MissLink{
		ProductIdFk:      productId,
		OcrProductNameFk: ocrName,
		UserId:           userId,
	}

	err := r.CreateEntity(&entity)
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *MissLinkRepository) IsLinkingDenied(productId int, ocrName string) (bool, error) {
	var missLinkCount int64 = 0
	err := r.Db.Model(&entities.MissLink{}).
		Where("product_id_fk = ? and ocr_product_name_fk = ?", productId, ocrName).
		Count(&missLinkCount).
		Error
	if err != nil {
		return true, err
	}

	return missLinkCount >= missLinkDenyLimit, nil
}

func (r *MissLinkRepository) ShouldBreakProductLink(productId int, ocrName string) (bool, error) {
	var linkCount int64 = -1
	err := r.Db.
		Table("ocr-product_product").
		Where(
			"product_entity_id = ? and ocr_product_entity_ocr_product_name= ?",
			productId,
			ocrName,
		).
		Count(&linkCount).
		Error
	if err != nil {
		return false, err
	}

	isLinkingDenied, err := r.IsLinkingDenied(productId, ocrName)
	if err != nil {
		return false, err
	}

	return linkCount > 0 && isLinkingDenied, nil
}

func (r *MissLinkRepository) GetReportsByUser(userId uint) ([]models.MissLink, error) {
	entitiesList := make([]entities.MissLink, 0)
	err := r.Db.Preload("Product").Preload("OcrProduct").Where("user_id = ?", userId).Find(&entitiesList).Error
	if err != nil {
		return nil, err
	}

	modelsList := make([]models.MissLink, 0)
	for _, entity := range entitiesList {
		model, err := toModel(entity)
		if err != nil {
			return nil, err
		}
		modelsList = append(modelsList, model)
	}

	return modelsList, nil
}

//func (r *MissLinkRepository) getMissLinksForOcrProducts(ocrNames []string) ([]entities.MissLink, error) {
//	var missLinks []entities.MissLink
//	err := r.Db.Where("ocr_product_name_fk IN (?)", ocrNames).Find(&missLinks).Error
//	if err != nil {
//		return nil, err
//	}
//
//	return missLinks, nil
//}

func (r *MissLinkRepository) getDeniedLinksForOcrProducts(ocrNames []string) (ocrProductsLinksDenied, error) {
	rows, err := r.Db.
		Table("miss_links").
		Select("ocr_product_name_fk, product_id_fk,  count(*) as count").
		Where("ocr_product_name_fk IN (?)", ocrNames).
		Group("ocr_product_name_fk, product_id_fk").
		Having("count(*) >= ?", missLinkDenyLimit).
		Rows()

	if err != nil {
		return nil, err
	}

	defer helpers.SafeClose(rows)

	deniedLinks := make(ocrProductsLinksDenied)
	for rows.Next() {
		var ocrName string
		var productId uint
		var count int
		err = rows.Scan(&ocrName, &productId, &count)
		if err != nil {
			return nil, err
		}
		deniedLinks[ocrNameAndProductId{ocrName, productId}] = struct{}{}
	}

	return deniedLinks, nil
}

func toModel(entity entities.MissLink) (models.MissLink, error) {
	return *models.NewMissLinkModelFromEntity(&entity), nil
}

func toEntity(model models.MissLink) (*entities.MissLink, error) {
	return model.ToEntity(), nil
}

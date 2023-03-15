package repository

import (
	"gorm.io/gorm/clause"
	"lib/data/database"
	"lib/data/database/repositories"
	"notifications/data/database/entities"
	"notifications/data/models"
)

type NotificationUserRepository struct {
	repositories.DbRepositoryWithModel[entities.NotificationUser, models.NotificationUserModel]
}

var repo *NotificationUserRepository = nil

func GetNotificationUserRepository() *NotificationUserRepository {
	if repo == nil {
		repo = &NotificationUserRepository{}
		db, err := database.GetDb()
		if err != nil {
			panic(err)
		}

		repo.Db = db

	}
	return repo
}

func (r *NotificationUserRepository) GetTokensByUserIds(userIds []int) ([]string, error) {
	var userEntities []entities.NotificationUser
	err := r.Db.Model(&entities.NotificationUser{}).
		Where("user_id IN ?", userIds).
		Find(&userEntities).
		Error
	if err != nil {
		return nil, err
	}
	tokens := make([]string, len(userEntities))
	for i, entity := range userEntities {
		tokens[i] = entity.FCMToken
	}
	return tokens, nil
}

func (r *NotificationUserRepository) CreateOrUpdate(userId int, fcmToken string) error {
	// TODO: Check sql
	return r.Db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"fcm_token"}),
	}).Create(&entities.NotificationUser{
		UserId:   userId,
		FCMToken: fcmToken,
	}).Error
}

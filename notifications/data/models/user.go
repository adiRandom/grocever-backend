package models

import "notifications/data/database/entities"

type NotificationUserModel struct {
	UserId   int    `json:"userId"`
	FCMToken string `json:"fcmToken"`
}

func NewFromNotificationUserEntity(entity *entities.NotificationUser) *NotificationUserModel {
	return &NotificationUserModel{
		UserId:   entity.UserId,
		FCMToken: entity.FCMToken,
	}
}

func (m *NotificationUserModel) ToEntity() *entities.NotificationUser {
	return &entities.NotificationUser{
		UserId:   m.UserId,
		FCMToken: m.FCMToken,
	}
}

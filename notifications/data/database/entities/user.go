package entities

import "gorm.io/gorm"

type NotificationUser struct {
	gorm.Model
	UserId   int `gorm:"uniqueIndex"`
	FCMToken string
}

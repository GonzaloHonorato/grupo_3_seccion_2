package entities

import "time"

type UserNotification struct {
	ID                     int        `json:"id"`
	UserID                 string     `json:"userId"`
	NotificationTemplateID int        `json:"notificationTemplateId"`
	IsRead                 bool       `json:"isRead"`
	ReadAt                 *time.Time `json:"readAt"` 
}

type UserNotifications []UserNotification

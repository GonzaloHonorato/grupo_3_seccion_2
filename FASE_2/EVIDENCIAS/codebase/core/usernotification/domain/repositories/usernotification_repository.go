package repositories

import "github.com/gonzalohonorato/servercorego/core/usernotification/domain/entities"

type UserNotificationRepository interface {
	SearchUserNotificationByID(id int) (*entities.UserNotification, error)
	SearchUserNotifications() (*entities.UserNotifications, error)
	SearchUserNotificationsByUserID(userID string) (*entities.UserNotifications, error)
	CreateUserNotification(notification *entities.UserNotification) error
	UpdateUserNotificationByID(notification *entities.UserNotification) error
	DeleteUserNotificationByID(id int) error
}

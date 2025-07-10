package application

import (
	"github.com/gonzalohonorato/servercorego/core/usernotification/domain/entities"
	"github.com/gonzalohonorato/servercorego/core/usernotification/domain/repositories"
)

type UserNotificationUsecase struct {
	UserNotificationRepository repositories.UserNotificationRepository
}

func NewUserNotificationUsecase(repo repositories.UserNotificationRepository) *UserNotificationUsecase {
	return &UserNotificationUsecase{UserNotificationRepository: repo}
}

func (uc *UserNotificationUsecase) SearchUserNotificationByID(id int) (*entities.UserNotification, error) {
	return uc.UserNotificationRepository.SearchUserNotificationByID(id)
}

func (uc *UserNotificationUsecase) SearchUserNotifications() (*entities.UserNotifications, error) {
	return uc.UserNotificationRepository.SearchUserNotifications()
}

func (uc *UserNotificationUsecase) CreateUserNotification(notification *entities.UserNotification) error {
	return uc.UserNotificationRepository.CreateUserNotification(notification)
}

func (uc *UserNotificationUsecase) UpdateUserNotificationByID(notification *entities.UserNotification) error {
	return uc.UserNotificationRepository.UpdateUserNotificationByID(notification)
}

func (uc *UserNotificationUsecase) DeleteUserNotificationByID(id int) error {
	return uc.UserNotificationRepository.DeleteUserNotificationByID(id)
}

func (uc *UserNotificationUsecase) SearchUserNotificationsByUserID(userID string) (*entities.UserNotifications, error) {
	return uc.UserNotificationRepository.SearchUserNotificationsByUserID(userID)
}

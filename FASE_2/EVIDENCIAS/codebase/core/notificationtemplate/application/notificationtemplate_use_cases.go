package application

import (
	"github.com/gonzalohonorato/servercorego/core/notificationtemplate/domain/entities"
	"github.com/gonzalohonorato/servercorego/core/notificationtemplate/domain/repositories"
)

type NotificationTemplateUsecase struct {
	NotificationTemplateRepository repositories.NotificationTemplateRepository
}

func NewNotificationTemplateUsecase(repo repositories.NotificationTemplateRepository) *NotificationTemplateUsecase {
	return &NotificationTemplateUsecase{NotificationTemplateRepository: repo}
}

func (uc *NotificationTemplateUsecase) SearchNotificationTemplateByID(id int) (*entities.NotificationTemplate, error) {
	return uc.NotificationTemplateRepository.SearchNotificationTemplateByID(id)
}

func (uc *NotificationTemplateUsecase) SearchNotificationTemplates() (*entities.NotificationTemplates, error) {
	return uc.NotificationTemplateRepository.SearchNotificationTemplates()
}

func (uc *NotificationTemplateUsecase) CreateNotificationTemplate(template *entities.NotificationTemplate) error {
	return uc.NotificationTemplateRepository.CreateNotificationTemplate(template)
}

func (uc *NotificationTemplateUsecase) UpdateNotificationTemplateByID(template *entities.NotificationTemplate) error {
	return uc.NotificationTemplateRepository.UpdateNotificationTemplateByID(template)
}

func (uc *NotificationTemplateUsecase) DeleteNotificationTemplateByID(id int) error {
	return uc.NotificationTemplateRepository.DeleteNotificationTemplateByID(id)
}

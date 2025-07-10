package repositories

import "github.com/gonzalohonorato/servercorego/core/notificationtemplate/domain/entities"

type NotificationTemplateRepository interface {
	SearchNotificationTemplateByID(id int) (*entities.NotificationTemplate, error)
	SearchNotificationTemplates() (*entities.NotificationTemplates, error)
	CreateNotificationTemplate(template *entities.NotificationTemplate) error
	UpdateNotificationTemplateByID(template *entities.NotificationTemplate) error
	DeleteNotificationTemplateByID(id int) error
}

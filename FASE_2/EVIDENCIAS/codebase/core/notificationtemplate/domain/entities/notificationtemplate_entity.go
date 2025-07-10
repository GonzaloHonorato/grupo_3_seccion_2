package entities

import "time"

type NotificationTemplate struct {
	ID        int        `json:"id"`
	Title     string     `json:"title"`
	Message   string     `json:"message"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
}

type NotificationTemplates []NotificationTemplate

package persistence

import (
	"context"

	"github.com/gonzalohonorato/servercorego/core/notificationtemplate/domain/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TimescaleNotificationTemplateRepository struct {
	dbPool *pgxpool.Pool
}

func NewTimescaleNotificationTemplateRepository(pool *pgxpool.Pool) *TimescaleNotificationTemplateRepository {
	return &TimescaleNotificationTemplateRepository{
		dbPool: pool,
	}
}

func (r *TimescaleNotificationTemplateRepository) SearchNotificationTemplateByID(id int) (*entities.NotificationTemplate, error) {
	ctx := context.Background()
	query := `SELECT id, title, message, created_at FROM notification_template WHERE id = $1`
	row := r.dbPool.QueryRow(ctx, query, id)
	var nt entities.NotificationTemplate
	err := row.Scan(&nt.ID, &nt.Title, &nt.Message, &nt.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &nt, nil
}

func (r *TimescaleNotificationTemplateRepository) SearchNotificationTemplates() (*entities.NotificationTemplates, error) {
	ctx := context.Background()
	query := `SELECT id, title, message, created_at FROM notification_template`
	rows, err := r.dbPool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates entities.NotificationTemplates
	for rows.Next() {
		var nt entities.NotificationTemplate
		if err := rows.Scan(&nt.ID, &nt.Title, &nt.Message, &nt.CreatedAt); err != nil {
			return nil, err
		}
		templates = append(templates, nt)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &templates, nil
}

func (r *TimescaleNotificationTemplateRepository) CreateNotificationTemplate(nt *entities.NotificationTemplate) error {
	ctx := context.Background()
	query := `
	INSERT INTO notification_template (
		title, message, created_at
	) VALUES (
		$1, $2, $3
	) RETURNING id
	`

	
	err := r.dbPool.QueryRow(ctx, query, nt.Title, nt.Message, nt.CreatedAt).Scan(&nt.ID)
	return err
}

func (r *TimescaleNotificationTemplateRepository) UpdateNotificationTemplateByID(nt *entities.NotificationTemplate) error {
	ctx := context.Background()
	query := `UPDATE notification_template SET title = $2, message = $3, created_at = $4 WHERE id = $1`
	_, err := r.dbPool.Exec(ctx, query, nt.ID, nt.Title, nt.Message, nt.CreatedAt)
	return err
}

func (r *TimescaleNotificationTemplateRepository) DeleteNotificationTemplateByID(id int) error {
	ctx := context.Background()
	query := `DELETE FROM notification_template WHERE id = $1`
	_, err := r.dbPool.Exec(ctx, query, id)
	return err
}

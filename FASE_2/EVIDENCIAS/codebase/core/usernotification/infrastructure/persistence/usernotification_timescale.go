package persistence

import (
	"context"
	"fmt"

	"github.com/gonzalohonorato/servercorego/core/usernotification/domain/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TimescaleUserNotificationRepository struct {
	dbPool *pgxpool.Pool
}

func NewTimescaleUserNotificationRepository(pool *pgxpool.Pool) *TimescaleUserNotificationRepository {
	return &TimescaleUserNotificationRepository{
		dbPool: pool,
	}
}

func (r *TimescaleUserNotificationRepository) SearchUserNotificationByID(id int) (*entities.UserNotification, error) {
	ctx := context.Background()
	query := `SELECT id, user_id, notification_template_id, is_read, read_at FROM user_notification WHERE id = $1`
	row := r.dbPool.QueryRow(ctx, query, id)
	var un entities.UserNotification
	err := row.Scan(&un.ID, &un.UserID, &un.NotificationTemplateID, &un.IsRead, &un.ReadAt)
	if err != nil {
		return nil, err
	}
	return &un, nil
}

func (r *TimescaleUserNotificationRepository) SearchUserNotifications() (*entities.UserNotifications, error) {
	ctx := context.Background()
	query := `SELECT id, user_id, notification_template_id, is_read, read_at FROM user_notification`
	rows, err := r.dbPool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications entities.UserNotifications
	for rows.Next() {
		var un entities.UserNotification
		if err := rows.Scan(&un.ID, &un.UserID, &un.NotificationTemplateID, &un.IsRead, &un.ReadAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, un)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &notifications, nil
}

func (r *TimescaleUserNotificationRepository) CreateUserNotification(un *entities.UserNotification) error {
	ctx := context.Background()
	query := `
	INSERT INTO user_notification (
		user_id, notification_template_id, is_read, read_at
	) VALUES (
		$1, $2, $3, $4
	);
	`

	
	fmt.Printf("Inserting UserNotification: UserID=%s, TemplateID=%d, IsRead=%t\n",
		un.UserID, un.NotificationTemplateID, un.IsRead)

	_, err := r.dbPool.Exec(ctx, query, un.UserID, un.NotificationTemplateID, un.IsRead, un.ReadAt)
	return err
}

func (r *TimescaleUserNotificationRepository) UpdateUserNotificationByID(un *entities.UserNotification) error {
	ctx := context.Background()
	query := `UPDATE user_notification SET user_id = $2, notification_template_id = $3, is_read = $4, read_at = $5 WHERE id = $1`
	_, err := r.dbPool.Exec(ctx, query, un.ID, un.UserID, un.NotificationTemplateID, un.IsRead, un.ReadAt)
	return err
}

func (r *TimescaleUserNotificationRepository) DeleteUserNotificationByID(id int) error {
	ctx := context.Background()
	query := `DELETE FROM user_notification WHERE id = $1`
	_, err := r.dbPool.Exec(ctx, query, id)
	return err
}

func (r *TimescaleUserNotificationRepository) SearchUserNotificationsByUserID(userID string) (*entities.UserNotifications, error) {
	ctx := context.Background()
	query := `SELECT id, user_id, notification_template_id, is_read, read_at FROM user_notification WHERE user_id = $1 ORDER BY id DESC`
	rows, err := r.dbPool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications entities.UserNotifications
	for rows.Next() {
		var un entities.UserNotification
		if err := rows.Scan(&un.ID, &un.UserID, &un.NotificationTemplateID, &un.IsRead, &un.ReadAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, un)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &notifications, nil
}

package persistence

import (
	"context"
	"time"

	"github.com/gonzalohonorato/servercorego/core/feedback/domain/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TimescaleFeedbackRepository struct {
	dbPool *pgxpool.Pool
}

func NewTimescaleDBRepository(pool *pgxpool.Pool) *TimescaleFeedbackRepository {
	return &TimescaleFeedbackRepository{
		dbPool: pool,
	}
}

func (r *TimescaleFeedbackRepository) SearchFeedbackByID(id int) (*entities.Feedback, error) {
	ctx := context.Background()
	query := `SELECT id, from_user_id, comment, response_comment, rating, created_at, response_at FROM feedback WHERE id = $1`
	row := r.dbPool.QueryRow(ctx, query, id)
	var f entities.Feedback
	err := row.Scan(&f.ID, &f.FromUserID, &f.Comment, &f.ResponseComment, &f.Rating, &f.CreatedAt, &f.ResponseAt)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *TimescaleFeedbackRepository) SearchFeedbacks() (*entities.Feedbacks, error) {
	ctx := context.Background()
	query := `SELECT id, from_user_id, comment, response_comment, rating, created_at, response_at FROM feedback`
	rows, err := r.dbPool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbacks entities.Feedbacks
	for rows.Next() {
		var f entities.Feedback
		if err := rows.Scan(&f.ID, &f.FromUserID, &f.Comment, &f.ResponseComment, &f.Rating, &f.CreatedAt, &f.ResponseAt); err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, f)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &feedbacks, nil
}

func (r *TimescaleFeedbackRepository) CreateFeedback(feedback *entities.Feedback) error {
	ctx := context.Background()
	query := `
	INSERT INTO feedback (
		from_user_id, comment, response_comment, rating, created_at, response_at
	) VALUES (
		$1, $2, $3, $4, $5, $6
	);
`
	
	if feedback.CreatedAt.IsZero() {
		feedback.CreatedAt = time.Now()
	}

	_, err := r.dbPool.Exec(ctx, query,
		feedback.FromUserID,
		feedback.Comment,
		feedback.ResponseComment,
		feedback.Rating,
		feedback.CreatedAt,
		feedback.ResponseAt)
	return err
}

func (r *TimescaleFeedbackRepository) UpdateFeedbackByID(f *entities.Feedback) error {
	ctx := context.Background()
	query := `UPDATE feedback SET from_user_id = $2, comment = $3, response_comment = $4, rating = $5, created_at = $6, response_at = $7 WHERE id = $1`

	_, err := r.dbPool.Exec(ctx, query,
		f.ID,
		f.FromUserID,
		f.Comment,
		f.ResponseComment,
		f.Rating,
		f.CreatedAt,
		f.ResponseAt)
	return err
}

func (r *TimescaleFeedbackRepository) DeleteFeedbackByID(id int) error {
	ctx := context.Background()
	query := `DELETE FROM feedback WHERE id = $1`
	_, err := r.dbPool.Exec(ctx, query, id)
	return err
}

func (r *TimescaleFeedbackRepository) SearchFeedbacksByUserID(userID string) (*entities.Feedbacks, error) {
	ctx := context.Background()
	query := `SELECT id, from_user_id, comment, response_comment, rating, created_at, response_at FROM feedback WHERE from_user_id = $1`
	rows, err := r.dbPool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbacks entities.Feedbacks
	for rows.Next() {
		var f entities.Feedback
		if err := rows.Scan(&f.ID, &f.FromUserID, &f.Comment, &f.ResponseComment, &f.Rating, &f.CreatedAt, &f.ResponseAt); err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, f)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &feedbacks, nil
}

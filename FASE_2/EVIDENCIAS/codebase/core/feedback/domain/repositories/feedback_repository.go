package repositories

import "github.com/gonzalohonorato/servercorego/core/feedback/domain/entities"

type FeedbackRepository interface {
	SearchFeedbackByID(id int) (*entities.Feedback, error)
	SearchFeedbacks() (*entities.Feedbacks, error)
	SearchFeedbacksByUserID(userID string) (*entities.Feedbacks, error)
	CreateFeedback(feedback *entities.Feedback) error
	UpdateFeedbackByID(feedback *entities.Feedback) error
	DeleteFeedbackByID(id int) error
}

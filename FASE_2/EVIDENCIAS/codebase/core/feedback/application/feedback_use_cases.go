package application

import (
	"github.com/gonzalohonorato/servercorego/core/feedback/domain/entities"
	"github.com/gonzalohonorato/servercorego/core/feedback/domain/repositories"
)

type FeedbackUsecase struct {
	FeedbackRepository repositories.FeedbackRepository
}

func NewFeedbackUsecase(feedbackRepo repositories.FeedbackRepository) *FeedbackUsecase {
	return &FeedbackUsecase{FeedbackRepository: feedbackRepo}
}

func (uc *FeedbackUsecase) SearchFeedbackByID(id int) (*entities.Feedback, error) {
	return uc.FeedbackRepository.SearchFeedbackByID(id)
}

func (uc *FeedbackUsecase) SearchFeedbacks() (*entities.Feedbacks, error) {
	return uc.FeedbackRepository.SearchFeedbacks()
}

func (uc *FeedbackUsecase) CreateFeedback(feedback *entities.Feedback) error {
	return uc.FeedbackRepository.CreateFeedback(feedback)
}

func (uc *FeedbackUsecase) UpdateFeedbackById(feedback *entities.Feedback) error {
	return uc.FeedbackRepository.UpdateFeedbackByID(feedback)
}

func (uc *FeedbackUsecase) DeleteFeedbackByID(id int) error {
	return uc.FeedbackRepository.DeleteFeedbackByID(id)
}

func (uc *FeedbackUsecase) SearchFeedbacksByUserID(userID string) (*entities.Feedbacks, error) {
	return uc.FeedbackRepository.SearchFeedbacksByUserID(userID)
}

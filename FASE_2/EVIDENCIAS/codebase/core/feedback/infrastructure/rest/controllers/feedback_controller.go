package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gonzalohonorato/servercorego/core/feedback/application"
	"github.com/gonzalohonorato/servercorego/core/feedback/domain/entities"
	"github.com/gonzalohonorato/servercorego/core/feedback/domain/repositories"
	"github.com/gorilla/mux"
)

type FeedbackController struct {
	FeedbackUsecase *application.FeedbackUsecase
}

func NewFeedbackController(feedbackRepository repositories.FeedbackRepository) *FeedbackController {
	feedbackUseCase := application.NewFeedbackUsecase(feedbackRepository)

	return &FeedbackController{
		FeedbackUsecase: feedbackUseCase,
	}
}

func (uc *FeedbackController) GetFeedbackByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	feedbackID := vars["id"]
	idInt, err := strconv.Atoi(feedbackID)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	feedback, err := uc.FeedbackUsecase.SearchFeedbackByID(idInt)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Feedback not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feedback)
}

func (uc *FeedbackController) GetFeedbacks(w http.ResponseWriter, r *http.Request) {
	feedbacks, err := uc.FeedbackUsecase.SearchFeedbacks()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Feedbacks not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feedbacks)
}

func (uc *FeedbackController) PostFeedback(w http.ResponseWriter, r *http.Request) {
	var newFeedback entities.Feedback
	if err := json.NewDecoder(r.Body).Decode(&newFeedback); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := uc.FeedbackUsecase.CreateFeedback(&newFeedback); err != nil {
		fmt.Println(err)
		http.Error(w, "Error creating feedback", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (uc *FeedbackController) PutFeedback(w http.ResponseWriter, r *http.Request) {
	var newFeedback entities.Feedback
	if err := json.NewDecoder(r.Body).Decode(&newFeedback); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := uc.FeedbackUsecase.UpdateFeedbackById(&newFeedback); err != nil {
		http.Error(w, "Error update feedback", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (uc *FeedbackController) DeleteFeedbackByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	feedbackID := vars["id"]
	idInt, err := strconv.Atoi(feedbackID)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	err = uc.FeedbackUsecase.DeleteFeedbackByID(idInt)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Feedback not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}

func (uc *FeedbackController) GetFeedbacksByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	feedbacks, err := uc.FeedbackUsecase.SearchFeedbacksByUserID(userID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error retrieving feedbacks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feedbacks)
}

func (uc *FeedbackController) RespondToFeedback(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	feedbackID := vars["id"]

	idInt, err := strconv.Atoi(feedbackID)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	
	var responseData struct {
		ResponseComment string `json:"responseComment"`
	}

	if err := json.NewDecoder(r.Body).Decode(&responseData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	
	feedback, err := uc.FeedbackUsecase.SearchFeedbackByID(idInt)
	if err != nil {
		http.Error(w, "Feedback not found", http.StatusNotFound)
		return
	}

	
	now := time.Now()
	feedback.ResponseComment = responseData.ResponseComment
	feedback.ResponseAt = &now

	if err := uc.FeedbackUsecase.UpdateFeedbackById(feedback); err != nil {
		http.Error(w, "Error updating feedback with response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

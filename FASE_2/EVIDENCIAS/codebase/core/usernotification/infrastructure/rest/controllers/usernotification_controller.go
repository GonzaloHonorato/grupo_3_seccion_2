package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gonzalohonorato/servercorego/core/usernotification/application"
	"github.com/gonzalohonorato/servercorego/core/usernotification/domain/entities"
	"github.com/gonzalohonorato/servercorego/core/usernotification/domain/repositories"
	"github.com/gorilla/mux"
)

type UserNotificationController struct {
	UserNotificationUsecase *application.UserNotificationUsecase
}

func NewUserNotificationController(userNotificationRepository repositories.UserNotificationRepository) *UserNotificationController {
	userNotificationUseCase := application.NewUserNotificationUsecase(userNotificationRepository)

	return &UserNotificationController{
		UserNotificationUsecase: userNotificationUseCase,
	}
}

func (uc *UserNotificationController) GetUserNotificationByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	notificationID := vars["id"]
	idInt, err := strconv.Atoi(notificationID)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	notification, err := uc.UserNotificationUsecase.SearchUserNotificationByID(idInt)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "User notification not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notification)
}

func (uc *UserNotificationController) GetUserNotifications(w http.ResponseWriter, r *http.Request) {
	notifications, err := uc.UserNotificationUsecase.SearchUserNotifications()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "User notifications not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

func (uc *UserNotificationController) PostUserNotification(w http.ResponseWriter, r *http.Request) {
	var newNotification entities.UserNotification
	if err := json.NewDecoder(r.Body).Decode(&newNotification); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := uc.UserNotificationUsecase.CreateUserNotification(&newNotification); err != nil {
		http.Error(w, "Error creating user notification", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (uc *UserNotificationController) PutUserNotification(w http.ResponseWriter, r *http.Request) {
	var newNotification entities.UserNotification
	if err := json.NewDecoder(r.Body).Decode(&newNotification); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := uc.UserNotificationUsecase.UpdateUserNotificationByID(&newNotification); err != nil {
		http.Error(w, "Error updating user notification", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (uc *UserNotificationController) DeleteUserNotificationByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	notificationID := vars["id"]
	idInt, err := strconv.Atoi(notificationID)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	err = uc.UserNotificationUsecase.DeleteUserNotificationByID(idInt)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "User notification not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}

func (uc *UserNotificationController) GetUserNotificationsByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	notifications, err := uc.UserNotificationUsecase.SearchUserNotificationsByUserID(userID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "User notifications not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

func (uc *UserNotificationController) MarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	notificationID := vars["id"]
	idInt, err := strconv.Atoi(notificationID)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	
	notification, err := uc.UserNotificationUsecase.SearchUserNotificationByID(idInt)
	if err != nil {
		http.Error(w, "User notification not found", http.StatusNotFound)
		return
	}

	
	if !notification.IsRead {
		now := time.Now()
		notification.IsRead = true
		notification.ReadAt = &now

		if err := uc.UserNotificationUsecase.UpdateUserNotificationByID(notification); err != nil {
			http.Error(w, "Error updating notification", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gonzalohonorato/servercorego/core/notificationtemplate/application"
	"github.com/gonzalohonorato/servercorego/core/notificationtemplate/domain/entities"
	"github.com/gonzalohonorato/servercorego/core/notificationtemplate/domain/repositories"
	userapp "github.com/gonzalohonorato/servercorego/core/user/application"
	userrepo "github.com/gonzalohonorato/servercorego/core/user/domain/repositories"
	usernotificationapp "github.com/gonzalohonorato/servercorego/core/usernotification/application"
	usernotificationentities "github.com/gonzalohonorato/servercorego/core/usernotification/domain/entities"
	usernotificationrepo "github.com/gonzalohonorato/servercorego/core/usernotification/domain/repositories"
	"github.com/gonzalohonorato/servercorego/core/websocket/infrastructure"
	"github.com/gorilla/mux"
)

type NotificationTemplateController struct {
	NotificationTemplateUsecase *application.NotificationTemplateUsecase
	UserNotificationUsecase     *usernotificationapp.UserNotificationUsecase
	UserUsecase                 *userapp.UserUsecase             
	WebSocketService            *infrastructure.WebSocketService 
}


func NewNotificationTemplateController(
	notificationTemplateRepository repositories.NotificationTemplateRepository,
	userNotificationRepository usernotificationrepo.UserNotificationRepository,
	userRepository userrepo.UserRepository,
	wsService *infrastructure.WebSocketService, 
) *NotificationTemplateController {
	
	notificationTemplateUsecase := application.NewNotificationTemplateUsecase(notificationTemplateRepository)
	userNotificationUsecase := usernotificationapp.NewUserNotificationUsecase(userNotificationRepository)
	userUsecase := userapp.NewUserUsecase(userRepository)

	return &NotificationTemplateController{
		NotificationTemplateUsecase: notificationTemplateUsecase,
		UserNotificationUsecase:     userNotificationUsecase,
		UserUsecase:                 userUsecase, 
		WebSocketService:            wsService,   
	}
}


type SendNotificationRequest struct {
	Title         string   `json:"title"`
	Message       string   `json:"message"`
	RecipientType string   `json:"recipientType"` 
	Groups        []string `json:"groups,omitempty"`
	UserIDs       []string `json:"userIds,omitempty"`
}


type SendNotificationResponse struct {
	Message       string `json:"message"`
	TemplateID    int    `json:"templateId"`
	TotalSent     int    `json:"totalSent"`
	TotalUsers    int    `json:"totalUsers"`
	RecipientType string `json:"recipientType"`
}

func (uc *NotificationTemplateController) GetNotificationTemplateByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateID := vars["id"]
	idInt, err := strconv.Atoi(templateID)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	template, err := uc.NotificationTemplateUsecase.SearchNotificationTemplateByID(idInt)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "NotificationTemplate not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(template)
}

func (uc *NotificationTemplateController) GetNotificationTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := uc.NotificationTemplateUsecase.SearchNotificationTemplates()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "NotificationTemplates not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(templates)
}

func (uc *NotificationTemplateController) PostNotificationTemplate(w http.ResponseWriter, r *http.Request) {
	var newTemplate entities.NotificationTemplate
	if err := json.NewDecoder(r.Body).Decode(&newTemplate); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	
	if newTemplate.CreatedAt == nil {
		now := time.Now()
		newTemplate.CreatedAt = &now
	}

	if err := uc.NotificationTemplateUsecase.CreateNotificationTemplate(&newTemplate); err != nil {
		http.Error(w, "Error creating NotificationTemplate", http.StatusInternalServerError)
		return
	}

	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTemplate)
}

func (uc *NotificationTemplateController) PutNotificationTemplate(w http.ResponseWriter, r *http.Request) {
	var newTemplate entities.NotificationTemplate
	if err := json.NewDecoder(r.Body).Decode(&newTemplate); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := uc.NotificationTemplateUsecase.UpdateNotificationTemplateByID(&newTemplate); err != nil {
		http.Error(w, "Error updating NotificationTemplate", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (uc *NotificationTemplateController) DeleteNotificationTemplateByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateID := vars["id"]
	idInt, err := strconv.Atoi(templateID)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	err = uc.NotificationTemplateUsecase.DeleteNotificationTemplateByID(idInt)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "NotificationTemplate not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}

func (uc *NotificationTemplateController) SendNotification(w http.ResponseWriter, r *http.Request) {
	var request SendNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	
	if request.Title == "" || request.Message == "" || request.RecipientType == "" {
		http.Error(w, "Title, message and recipientType are required", http.StatusBadRequest)
		return
	}

	
	now := time.Now()
	template := &entities.NotificationTemplate{
		Title:     request.Title,
		Message:   request.Message,
		CreatedAt: &now,
	}

	if err := uc.NotificationTemplateUsecase.CreateNotificationTemplate(template); err != nil {
		http.Error(w, "Error creating notification template", http.StatusInternalServerError)
		return
	}

	fmt.Printf("✅ Template creado con ID: %d\n", template.ID)

	
	userIDs, err := uc.getUserIDsByRecipientType(request.RecipientType, request.Groups, request.UserIDs)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting recipient users: %v", err), http.StatusInternalServerError)
		return
	}

	if len(userIDs) == 0 {
		http.Error(w, "No users found for the specified criteria", http.StatusBadRequest)
		return
	}

	
	sentCount := 0
	notifiedUsers := make([]string, 0) 

	for _, userID := range userIDs {
		userNotification := &usernotificationentities.UserNotification{
			UserID:                 userID,
			NotificationTemplateID: template.ID,
			IsRead:                 false,
		}

		if err := uc.UserNotificationUsecase.CreateUserNotification(userNotification); err != nil {
			fmt.Printf("❌ Error creando notificación para usuario %s: %v\n", userID, err)
			continue
		}

		sentCount++
		notifiedUsers = append(notifiedUsers, userID)
		fmt.Printf("✅ Notificación creada para usuario %s con template ID: %d\n", userID, template.ID)
	}

	
	if len(notifiedUsers) > 0 && uc.WebSocketService != nil {
		
		notificationPayload := map[string]interface{}{
			"id":                     template.ID, 
			"notificationTemplateId": template.ID,
			"userId":                 "", 
			"template": map[string]interface{}{
				"id":        template.ID,
				"title":     template.Title,
				"message":   template.Message,
				"createdAt": template.CreatedAt.Format(time.RFC3339),
			},
			"isRead": false,
			"readAt": nil,
		}

		
		uc.WebSocketService.NotifyMultipleUsers(notifiedUsers, "new_notification", notificationPayload)

		fmt.Printf("🔔 WebSocket: Notificaciones enviadas a %d usuario(s)\n", len(notifiedUsers))

		
		stats := uc.WebSocketService.GetConnectionStats()
		fmt.Printf("📊 Conexiones WebSocket activas: %d admin, %d usuarios\n",
			stats["admin_clients"], stats["user_clients"])
	} else {
		fmt.Printf("⚠️ WebSocket service no disponible o sin usuarios para notificar\n")
	}

	
	response := SendNotificationResponse{
		Message:       "Notifications sent successfully",
		TemplateID:    template.ID,
		TotalSent:     sentCount,
		TotalUsers:    len(userIDs),
		RecipientType: request.RecipientType,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}


func (uc *NotificationTemplateController) getUserIDsByRecipientType(recipientType string, groups []string, userIDs []string) ([]string, error) {
	switch recipientType {
	case "all":
		return uc.getAllUserIDs()
	case "groups":
		return uc.getUserIDsByGroups(groups)
	case "specific":
		return userIDs, nil
	default:
		return nil, fmt.Errorf("invalid recipient type: %s", recipientType)
	}
}

func (uc *NotificationTemplateController) getAllUserIDs() ([]string, error) {
	users, err := uc.UserUsecase.SearchUsers() 
	if err != nil {
		return nil, err
	}

	var userIDs []string
	for _, user := range *users {
		userIDs = append(userIDs, user.ID)
	}
	return userIDs, nil
}

func (uc *NotificationTemplateController) getUserIDsByGroups(groups []string) ([]string, error) {
	var allUserIDs []string

	for _, group := range groups {
		users, err := uc.UserUsecase.SearchUsersByType(group) 
		if err != nil {
			return nil, err
		}

		for _, user := range *users {
			
			found := false
			for _, existingID := range allUserIDs {
				if existingID == user.ID {
					found = true
					break
				}
			}
			if !found {
				allUserIDs = append(allUserIDs, user.ID)
			}
		}
	}

	return allUserIDs, nil
}

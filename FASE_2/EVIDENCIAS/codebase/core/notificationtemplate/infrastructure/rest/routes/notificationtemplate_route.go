package routes

import (
	"github.com/gonzalohonorato/servercorego/config/injector"
	"github.com/gonzalohonorato/servercorego/core/notificationtemplate/infrastructure/rest/controllers"
	"github.com/gorilla/mux"
)

func NotificationTemplateRoutes(router *mux.Router, container *injector.Container) {
	controller := controllers.NewNotificationTemplateController(
		container.ProvideNotificationTemplateRepository(),
		container.ProvideUserNotificationRepository(),
		container.ProvideUserRepository(),
		container.ProvideWebSocketService(), 
	)

	router.HandleFunc("/notification-templates", controller.PutNotificationTemplate).Methods("PUT")
	router.HandleFunc("/notification-templates", controller.GetNotificationTemplates).Methods("GET")
	router.HandleFunc("/notification-templates/{id}", controller.DeleteNotificationTemplateByID).Methods("DELETE")
	router.HandleFunc("/notification-templates/{id}", controller.GetNotificationTemplateByID).Methods("GET")
	router.HandleFunc("/notification-templates", controller.PostNotificationTemplate).Methods("POST")
	router.HandleFunc("/notifications/send", controller.SendNotification).Methods("POST")
}

package routes

import (
	"github.com/gonzalohonorato/servercorego/config/injector"
	"github.com/gonzalohonorato/servercorego/core/usernotification/infrastructure/rest/controllers"
	"github.com/gorilla/mux"
)

func UserNotificationRoutes(router *mux.Router, container *injector.Container) {
	controller := controllers.NewUserNotificationController(container.ProvideUserNotificationRepository())
	router.HandleFunc("/user-notifications", controller.PutUserNotification).Methods("PUT")
	router.HandleFunc("/user-notifications", controller.GetUserNotifications).Methods("GET")
	router.HandleFunc("/user-notifications/{id}", controller.DeleteUserNotificationByID).Methods("DELETE")
	router.HandleFunc("/user-notifications/{id}", controller.GetUserNotificationByID).Methods("GET")
	router.HandleFunc("/user-notifications", controller.PostUserNotification).Methods("POST")
	router.HandleFunc("/users/{userId}/notifications", controller.GetUserNotificationsByUserID).Methods("GET")
	router.HandleFunc("/user-notifications/{id}/read", controller.MarkNotificationAsRead).Methods("POST")
}

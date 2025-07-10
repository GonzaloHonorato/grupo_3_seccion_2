package routes

import (
	"github.com/gonzalohonorato/servercorego/config/injector"
	"github.com/gonzalohonorato/servercorego/core/feedback/infrastructure/rest/controllers"
	"github.com/gorilla/mux"
)

func FeedbackRoutes(router *mux.Router, container *injector.Container) {
	controller := controllers.NewFeedbackController(container.ProvideFeedbackRepository())
	router.HandleFunc("/feedbacks", controller.PutFeedback).Methods("PUT")
	router.HandleFunc("/feedbacks", controller.GetFeedbacks).Methods("GET")
	router.HandleFunc("/feedbacks/{id}", controller.DeleteFeedbackByID).Methods("DELETE")
	router.HandleFunc("/feedbacks/{id}", controller.GetFeedbackByID).Methods("GET")
	router.HandleFunc("/feedbacks", controller.PostFeedback).Methods("POST")
	router.HandleFunc("/users/{userID}/feedbacks", controller.GetFeedbacksByUserID).Methods("GET")
	router.HandleFunc("/feedbacks/{id}/respond", controller.RespondToFeedback).Methods("POST")
}

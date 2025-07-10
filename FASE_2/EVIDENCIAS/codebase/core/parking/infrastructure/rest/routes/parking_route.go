package routes

import (
	"github.com/gonzalohonorato/servercorego/config/injector"
	"github.com/gonzalohonorato/servercorego/core/parking/infrastructure/rest/controllers"
	"github.com/gorilla/mux"
)

func ParkingRoutes(router *mux.Router, container *injector.Container) {
	controller := controllers.NewParkingController(container.ProvideParkingRepository())
	router.HandleFunc("/parkings", controller.PutParking).Methods("PUT")
	router.HandleFunc("/parkings", controller.GetParkings).Methods("GET")
	router.HandleFunc("/available-parkings", controller.GetAvailableParkings).Methods("GET")

	router.HandleFunc("/parkings/{id}", controller.DeleteParkingByID).Methods("DELETE")
	router.HandleFunc("/parkings/{id}", controller.GetParkingByID).Methods("GET")
	router.HandleFunc("/parkings", controller.PostParking).Methods("POST")
}

package routes

import (
	"github.com/gonzalohonorato/servercorego/config/injector"
	"github.com/gonzalohonorato/servercorego/core/vehicle/infrastructure/rest/controllers"
	"github.com/gorilla/mux"
)

func VehicleRoutes(router *mux.Router, container *injector.Container) {
	controller := controllers.NewVehicleController(container.ProvideVehicleRepository())
	router.HandleFunc("/vehicles", controller.PutVehicle).Methods("PUT")
	router.HandleFunc("/vehicles", controller.GetVehicles).Methods("GET")
	router.HandleFunc("/vehicles-user/{customerId}", controller.GetVehiclesByCustomerID).Methods("GET")
	router.HandleFunc("/vehicles/{id}", controller.DeleteVehicleByID).Methods("DELETE")
	router.HandleFunc("/vehicles/{id}", controller.GetVehicleByID).Methods("GET")
	router.HandleFunc("/vehicles", controller.PostVehicle).Methods("POST")
}

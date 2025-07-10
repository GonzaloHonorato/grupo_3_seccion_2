package routes

import (
	"github.com/gonzalohonorato/servercorego/config/injector"
	"github.com/gonzalohonorato/servercorego/core/parkingusage/infrastructure/rest/controllers"
	"github.com/gorilla/mux"
)

func ParkingUsageRoutes(router *mux.Router, container *injector.Container) {
	controller := controllers.NewParkingUsageController(
		container.ProvideParkingUsageRepository(),
		container.ProvideParkingRepository(),
		container.ProvideVehicleRepository(),
		container.ProvideReservationRepository(),
		container.ProvideWebSocketService(),
	)

	
	router.HandleFunc("/parking-usages", controller.PutParkingUsage).Methods("PUT")
	router.HandleFunc("/parking-usages", controller.GetParkingUsages).Methods("GET")
	router.HandleFunc("/active-parking-usages", controller.GetActiveParkingUsages).Methods("GET")
	router.HandleFunc("/parking-usages/vehicle/{id}", controller.GetParkingUsagesByVehicleID).Methods("GET")
	router.HandleFunc("/parking-usages/customer/{customerID}", controller.GetParkingUsagesByCustomerID).Methods("GET")

	router.HandleFunc("/parking-usages/{id}", controller.DeleteParkingUsageByID).Methods("DELETE")
	router.HandleFunc("/parking-usages/{id}", controller.GetParkingUsageByID).Methods("GET")
	router.HandleFunc("/parking-usages", controller.PostParkingUsage).Methods("POST")
	router.HandleFunc("/parking-usages/{id}/exit", controller.PostRegisterExitTime).Methods("POST")

	
	router.HandleFunc("/parking-entries", controller.PostParkingEntry).Methods("POST")
	router.HandleFunc("/parking-exits", controller.PostParkingExit).Methods("POST")
	router.HandleFunc("/ocr-entry", controller.PostOCREntry).Methods("POST")

	
	router.HandleFunc("/parking-entries/ocr", controller.PostOCREntryWithImage).Methods("POST")
	router.HandleFunc("/ocr/exits", controller.PostOCRExit).Methods("POST")
}

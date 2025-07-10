package routes

import (
	"github.com/gonzalohonorato/servercorego/config/injector"
	"github.com/gonzalohonorato/servercorego/core/reservation/infrastructure/rest/controllers"
	"github.com/gorilla/mux"
)

func ReservationRoutes(router *mux.Router, container *injector.Container) {
	controller := controllers.NewReservationController(
		container.ProvideReservationRepository(),
		container.ProvideParkingRepository(), 
	)

	router.HandleFunc("/reservations", controller.PutReservation).Methods("PUT")
	router.HandleFunc("/reservations", controller.GetReservations).Methods("GET")
	router.HandleFunc("/reservations-filter", controller.GetReservationsByDateAndStatus).Methods("GET")
	router.HandleFunc("/reservations-user", controller.GetReservationByUserIDAndDates).Methods("GET")
	router.HandleFunc("/reservations/{id}", controller.DeleteReservationByID).Methods("DELETE")
	router.HandleFunc("/reservations/{id}", controller.GetReservationByID).Methods("GET")
	router.HandleFunc("/reservations", controller.PostReservation).Methods("POST")
	router.HandleFunc("/reservations/{id}/status", controller.UpdateReservationStatus).Methods("PUT")
}

package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	parkingRepositories "github.com/gonzalohonorato/servercorego/core/parking/domain/repositories"
	"github.com/gonzalohonorato/servercorego/core/reservation/application"
	"github.com/gonzalohonorato/servercorego/core/reservation/domain/entities"
	"github.com/gonzalohonorato/servercorego/core/reservation/domain/repositories"
	"github.com/gorilla/mux"
)

type ReservationController struct {
	ReservationUsecase *application.ReservationUsecase
}

func NewReservationController(
	reservationRepository repositories.ReservationRepository,
	parkingRepository parkingRepositories.ParkingRepository, 
) *ReservationController {
	reservationUseCase := application.NewReservationUsecase(
		reservationRepository,
		parkingRepository, 
	)

	return &ReservationController{
		ReservationUsecase: reservationUseCase,
	}
}

func (uc *ReservationController) GetReservationByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reservationID := vars["id"]
	idInt, err := strconv.Atoi(reservationID)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	reservation, err := uc.ReservationUsecase.SearchReservationByID(idInt)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Reservation not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservation)
}

func (uc *ReservationController) GetReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := uc.ReservationUsecase.SearchReservations()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Reservations not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservations)
}

func (uc *ReservationController) GetReservationsByDateAndStatus(w http.ResponseWriter, r *http.Request) {
	
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		http.Error(w, "Date parameter is required", http.StatusBadRequest)
		return
	}

	
	statusValues := r.URL.Query()["status"]
	if len(statusValues) == 0 {
		http.Error(w, "At least one status parameter is required", http.StatusBadRequest)
		return
	}

	
	reservations, err := uc.ReservationUsecase.SearchReservationsByDateAndStatus(dateStr, statusValues)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error searching reservations: "+err.Error(), http.StatusInternalServerError)
		return
	}

	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservations)
}
func (uc *ReservationController) GetReservationByUserIDAndDates(w http.ResponseWriter, r *http.Request) {
	
	userID := r.URL.Query().Get("userID")
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	
	if userID == "" || startDate == "" || endDate == "" {
		http.Error(w, "userID, startDate, and endDate query parameters are required", http.StatusBadRequest)
		return
	}

	reservation, err := uc.ReservationUsecase.SearchReservationByUserIDAndDates(userID, startDate, endDate)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Reservation not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservation)
}
func (uc *ReservationController) PostReservation(w http.ResponseWriter, r *http.Request) {
	var newReservation entities.Reservation
	if err := json.NewDecoder(r.Body).Decode(&newReservation); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := uc.ReservationUsecase.CreateReservation(&newReservation); err != nil {
		http.Error(w, "Error creating reservation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (uc *ReservationController) PutReservation(w http.ResponseWriter, r *http.Request) {
	var newReservation entities.Reservation
	if err := json.NewDecoder(r.Body).Decode(&newReservation); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := uc.ReservationUsecase.UpdateReservationById(&newReservation); err != nil {
		http.Error(w, "Error update reservation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (uc *ReservationController) DeleteReservationByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reservationID := vars["id"]
	idInt, err := strconv.Atoi(reservationID)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	err = uc.ReservationUsecase.DeleteReservationByID(idInt)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Reservation not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}


func (uc *ReservationController) UpdateReservationStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID de reserva inválido", http.StatusBadRequest)
		return
	}

	var request struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	if err := uc.ReservationUsecase.UpdateReservationStatus(id, request.Status); err != nil {
		http.Error(w, fmt.Sprintf("Error al actualizar estado: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

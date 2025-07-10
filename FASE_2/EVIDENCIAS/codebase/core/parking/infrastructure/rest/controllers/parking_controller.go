package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gonzalohonorato/servercorego/core/parking/application"
	"github.com/gonzalohonorato/servercorego/core/parking/domain/entities"
	"github.com/gonzalohonorato/servercorego/core/parking/domain/repositories"
	"github.com/gorilla/mux"
)

type ParkingController struct {
	ParkingUsecase *application.ParkingUsecase
}

func NewParkingController(parkingRepository repositories.ParkingRepository) *ParkingController {
	parkingUseCase := application.NewParkingUsecase(parkingRepository)

	return &ParkingController{
		ParkingUsecase: parkingUseCase,
	}
}
func (uc *ParkingController) GetParkingByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	parkingID := vars["id"]
	idInt, err := strconv.Atoi(parkingID)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	parking, err := uc.ParkingUsecase.SearchParkingByID(idInt)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Parking not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(parking)
}

func (uc *ParkingController) GetAvailableParkings(w http.ResponseWriter, r *http.Request) {
	parkings, err := uc.ParkingUsecase.SearchAvailableParkings()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "No se pudieron encontrar estacionamientos disponibles", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(parkings)
}
func (uc *ParkingController) GetParkings(w http.ResponseWriter, r *http.Request) {
	parkings, err := uc.ParkingUsecase.SearchParkings()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Parkings not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(parkings)
}
func (uc *ParkingController) PostParking(w http.ResponseWriter, r *http.Request) {
	var newParking entities.Parking
	if err := json.NewDecoder(r.Body).Decode(&newParking); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := uc.ParkingUsecase.CreateParking(&newParking); err != nil {
		http.Error(w, "Error creating parking", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
func (uc *ParkingController) PutParking(w http.ResponseWriter, r *http.Request) {
	var newParking entities.Parking
	if err := json.NewDecoder(r.Body).Decode(&newParking); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := uc.ParkingUsecase.UpdateParkingById(&newParking); err != nil {
		http.Error(w, "Error update parking", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (uc *ParkingController) DeleteParkingByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	parkingID := vars["id"]
	idInt, err := strconv.Atoi(parkingID)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	err = uc.ParkingUsecase.DeleteParkingByID(idInt)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Parking not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}

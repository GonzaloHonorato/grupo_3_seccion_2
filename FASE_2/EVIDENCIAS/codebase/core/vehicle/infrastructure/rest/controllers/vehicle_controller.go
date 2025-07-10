package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gonzalohonorato/servercorego/core/vehicle/application"
	"github.com/gonzalohonorato/servercorego/core/vehicle/domain/entities"
	"github.com/gonzalohonorato/servercorego/core/vehicle/domain/repositories"
	"github.com/gorilla/mux"
)

type VehicleController struct {
	VehicleUsecase *application.VehicleUsecase
}

func NewVehicleController(vehicleRepository repositories.VehicleRepository) *VehicleController {
	vehicleUseCase := application.NewVehicleUsecase(vehicleRepository)

	return &VehicleController{
		VehicleUsecase: vehicleUseCase,
	}
}

func (uc *VehicleController) GetVehicleByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vehicleID := vars["id"]
	idInt, err := strconv.Atoi(vehicleID)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	vehicle, err := uc.VehicleUsecase.SearchVehicleByID(idInt)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Vehicle not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vehicle)
}

func (uc *VehicleController) GetVehiclesByCustomerID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID := vars["customerId"]
	vehicles, err := uc.VehicleUsecase.SearchVehiclesByCustomerID(customerID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Vehicles not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vehicles)
}

func (uc *VehicleController) GetVehicles(w http.ResponseWriter, r *http.Request) {
	vehicles, err := uc.VehicleUsecase.SearchVehicles()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Vehicles not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vehicles)
}

func (uc *VehicleController) PostVehicle(w http.ResponseWriter, r *http.Request) {
	var newVehicle entities.Vehicle
	if err := json.NewDecoder(r.Body).Decode(&newVehicle); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := uc.VehicleUsecase.CreateVehicle(&newVehicle); err != nil {
		http.Error(w, "Error creating vehicle", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (uc *VehicleController) PutVehicle(w http.ResponseWriter, r *http.Request) {
	var newVehicle entities.Vehicle
	if err := json.NewDecoder(r.Body).Decode(&newVehicle); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := uc.VehicleUsecase.UpdateVehicleById(&newVehicle); err != nil {
		http.Error(w, "Error update vehicle", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (uc *VehicleController) DeleteVehicleByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vehicleID := vars["id"]
	idInt, err := strconv.Atoi(vehicleID)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	err = uc.VehicleUsecase.DeleteVehicleByID(idInt)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Vehicle not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}

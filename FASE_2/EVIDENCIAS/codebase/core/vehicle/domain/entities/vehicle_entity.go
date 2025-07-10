package entities

import "time"

type Vehicle struct {
	ID          int       `json:"id"`
	Plate       string    `json:"plate"`
	Brand       string    `json:"brand"`
	Model       string    `json:"model"`
	VehicleType string    `json:"vehicleType"`
	CustomerID  string    `json:"customerId"` 
	CreatedAt   time.Time `json:"createdAt"`
}

type Vehicles []Vehicle

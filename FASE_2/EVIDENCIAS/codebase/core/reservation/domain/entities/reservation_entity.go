package entities

import "time"

type Reservation struct {
	ID         int       `json:"id"`
	CustomerID string    `json:"customerId"`
	ParkingID  int       `json:"parkingId"`
	VehicleID  int       `json:"vehicleId"`
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
	Status     string    `json:"status"` 
	CreatedAt  time.Time `json:"createdAt"`
}

type Reservations []Reservation

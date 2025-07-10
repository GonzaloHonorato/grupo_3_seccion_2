package entities

import "time"

type ParkingUsage struct {
	ID             int        `json:"id"`
	ReservationID  *int       `json:"reservationId"`
	VehicleID      *int       `json:"vehicleId"`
	ParkingID      int        `json:"parkingId"`
	EntryTime      *time.Time `json:"entryTime"`
	ExitTime       *time.Time `json:"exitTime"`
	OcrPlate       string     `json:"ocrPlate"`
	QrScanned      bool       `json:"qrScanned"`
	RegisteredBy   *string    `json:"registeredBy"`
	ManualEntry    bool       `json:"manualEntry"`
	VisitorName    string     `json:"visitorName"`
	VisitorRut     string     `json:"visitorRut"`
	VisitorContact string     `json:"visitorContact"`
	Zone           string     `json:"zone"`
}

type ParkingUsages []ParkingUsage

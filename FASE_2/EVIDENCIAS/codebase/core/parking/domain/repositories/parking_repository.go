package repositories

import "github.com/gonzalohonorato/servercorego/core/parking/domain/entities"

type ParkingRepository interface {
	SearchParkingByID(id int) (*entities.Parking, error)
	SearchParkings() (*entities.Parkings, error)
	SearchAvailableParkings() (*entities.Parkings, error)
	CreateParking(parking *entities.Parking) error
	UpdateParkingByID(parking *entities.Parking) error
	DeleteParkingByID(id int) error
}

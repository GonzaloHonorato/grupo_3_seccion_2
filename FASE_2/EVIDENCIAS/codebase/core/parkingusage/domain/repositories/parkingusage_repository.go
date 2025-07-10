package repositories

import "github.com/gonzalohonorato/servercorego/core/parkingusage/domain/entities"

type ParkingUsageRepository interface {
	SearchParkingUsageByID(id int) (*entities.ParkingUsage, error)
	SearchParkingUsages() (*entities.ParkingUsages, error)
	SearchParkingUsagesByVehicleID(id int) (*entities.ParkingUsages, error)
	SearchActiveParkingUsages() (*entities.ParkingUsages, error)
	CreateParkingUsage(parkingUsage *entities.ParkingUsage) error
	UpdateParkingUsageByID(parkingUsage *entities.ParkingUsage) error
	DeleteParkingUsageByID(id int) error
	SearchParkingUsagesByVehicleIDs(vehicleIDs []int, filters map[string]interface{}) (*[]entities.ParkingUsage, error)
}

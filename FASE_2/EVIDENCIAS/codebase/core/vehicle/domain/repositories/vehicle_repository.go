package repositories

import "github.com/gonzalohonorato/servercorego/core/vehicle/domain/entities"

type VehicleRepository interface {
	SearchVehicleByID(id int) (*entities.Vehicle, error)
	SearchVehiclesByCustomerID(customerID string) (*entities.Vehicles, error)
	SearchVehicles() (*entities.Vehicles, error)
	SearchVehicleByPlate(plate string) (*entities.Vehicle, error)
	CreateVehicle(vehicle *entities.Vehicle) error
	UpdateVehicleByID(vehicle *entities.Vehicle) error
	DeleteVehicleByID(id int) error
}

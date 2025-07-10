package application

import (
	"github.com/gonzalohonorato/servercorego/core/vehicle/domain/entities"
	"github.com/gonzalohonorato/servercorego/core/vehicle/domain/repositories"
)

type VehicleUsecase struct {
	VehicleRepository repositories.VehicleRepository
}

func NewVehicleUsecase(vehicleRepo repositories.VehicleRepository) *VehicleUsecase {
	return &VehicleUsecase{VehicleRepository: vehicleRepo}
}

func (uc *VehicleUsecase) SearchVehicleByID(id int) (*entities.Vehicle, error) {
	return uc.VehicleRepository.SearchVehicleByID(id)
}

func (uc *VehicleUsecase) SearchVehicles() (*entities.Vehicles, error) {
	return uc.VehicleRepository.SearchVehicles()
}
func (uc *VehicleUsecase) SearchVehiclesByCustomerID(customerID string) (*entities.Vehicles, error) {
	return uc.VehicleRepository.SearchVehiclesByCustomerID(customerID)
}
func (uc *VehicleUsecase) CreateVehicle(vehicle *entities.Vehicle) error {
	return uc.VehicleRepository.CreateVehicle(vehicle)
}

func (uc *VehicleUsecase) UpdateVehicleById(vehicle *entities.Vehicle) error {

	return uc.VehicleRepository.UpdateVehicleByID(vehicle)
}

func (uc *VehicleUsecase) DeleteVehicleByID(id int) error {
	return uc.VehicleRepository.DeleteVehicleByID(id)
}
func (uc *VehicleUsecase) SearchVehicleByPlate(plate string) (*entities.Vehicle, error) {
	return uc.VehicleRepository.SearchVehicleByPlate(plate)
}

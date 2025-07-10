package application

import (
	"github.com/gonzalohonorato/servercorego/core/parking/domain/entities"
	"github.com/gonzalohonorato/servercorego/core/parking/domain/repositories"
)

type ParkingUsecase struct {
	ParkingRepository repositories.ParkingRepository
}

func NewParkingUsecase(parkingRepo repositories.ParkingRepository) *ParkingUsecase {
	return &ParkingUsecase{ParkingRepository: parkingRepo}
}

func (uc *ParkingUsecase) SearchParkingByID(id int) (*entities.Parking, error) {
	return uc.ParkingRepository.SearchParkingByID(id)
}

func (uc *ParkingUsecase) SearchParkings() (*entities.Parkings, error) {
	return uc.ParkingRepository.SearchParkings()
}
func (uc *ParkingUsecase) SearchAvailableParkings() (*entities.Parkings, error) {
	return uc.ParkingRepository.SearchAvailableParkings()
}

func (uc *ParkingUsecase) CreateParking(parking *entities.Parking) error {

	return uc.ParkingRepository.CreateParking(parking)
}

func (uc *ParkingUsecase) UpdateParkingById(parking *entities.Parking) error {
	return uc.ParkingRepository.UpdateParkingByID(parking)
}

func (uc *ParkingUsecase) DeleteParkingByID(id int) error {
	return uc.ParkingRepository.DeleteParkingByID(id)
}

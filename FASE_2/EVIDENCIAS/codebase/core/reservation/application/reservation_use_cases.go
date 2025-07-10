package application

import (
	"fmt"
	"log"
	"time"

	parkingEntities "github.com/gonzalohonorato/servercorego/core/parking/domain/entities"
	parkingRepositories "github.com/gonzalohonorato/servercorego/core/parking/domain/repositories"
	"github.com/gonzalohonorato/servercorego/core/reservation/domain/entities"
	"github.com/gonzalohonorato/servercorego/core/reservation/domain/repositories"
)

type ReservationUsecase struct {
	ReservationRepository repositories.ReservationRepository
	ParkingRepository     parkingRepositories.ParkingRepository
}


func NewReservationUsecase(reservationRepo repositories.ReservationRepository, parkingRepo parkingRepositories.ParkingRepository) *ReservationUsecase {
	return &ReservationUsecase{
		ReservationRepository: reservationRepo,
		ParkingRepository:     parkingRepo,
	}
}

func (uc *ReservationUsecase) SearchReservationByID(id int) (*entities.Reservation, error) {
	return uc.ReservationRepository.SearchReservationByID(id)
}

func (uc *ReservationUsecase) SearchReservations() (*entities.Reservations, error) {
	return uc.ReservationRepository.SearchReservations()
}
func (uc *ReservationUsecase) SearchReservationsByDateAndStatus(date string, statuses []string) (*entities.Reservations, error) {
	return uc.ReservationRepository.SearchReservationsByDateAndStatus(date, statuses)
}

func (uc *ReservationUsecase) SearchReservationByUserIDAndDates(userID string, startDate string, endDate string) (*entities.Reservations, error) {
	return uc.ReservationRepository.SearchReservationByUserIDAndDates(userID, startDate, endDate)
}
func (uc *ReservationUsecase) CreateReservation(reservation *entities.Reservation) error {
	now := time.Now()
	timeUntilStart := reservation.StartTime.Sub(now)
	isImmediate := timeUntilStart <= time.Hour && timeUntilStart > 0

	
	if timeUntilStart < 0 {
		return fmt.Errorf("la fecha de inicio no puede ser en el pasado")
	}

	
	if reservation.ParkingID > 0 {
		
		ok, err := uc.isParkingAvailableForReservation(reservation.ParkingID, reservation.StartTime, reservation.EndTime, isImmediate)
		if err != nil {
			return err
		}
		if !ok {
			
			alt, err := uc.findAvailableParkingForReservation(reservation.StartTime, reservation.EndTime, isImmediate)
			if err != nil {
				return fmt.Errorf("no hay otro parking disponible: %w", err)
			}
			reservation.ParkingID = alt.ID
		}
	} else {
		
		p, err := uc.findAvailableParkingForReservation(reservation.StartTime, reservation.EndTime, isImmediate)
		if err != nil {
			return fmt.Errorf("no hay estacionamientos disponibles: %w", err)
		}
		reservation.ParkingID = p.ID
	}

	
	if isImmediate {
		reservation.Status = "active"
		if err := uc.activateParkingForReservation(reservation.ParkingID); err != nil {
			return fmt.Errorf("error al activar parking inmediato: %w", err)
		}
	} else {
		reservation.Status = "pending"
	}

	
	return uc.ReservationRepository.CreateReservation(reservation)
}


func (uc *ReservationUsecase) activateParkingForReservation(parkingID int) error {
	parking, err := uc.ParkingRepository.SearchParkingByID(parkingID)
	if err != nil {
		return fmt.Errorf("error al obtener parking: %w", err)
	}

	
	parking.IsActive = true
	err = uc.ParkingRepository.UpdateParkingByID(parking)
	if err != nil {
		return fmt.Errorf("error al activar parking: %w", err)
	}

	return nil
}
func (uc *ReservationUsecase) UpdateReservationById(reservation *entities.Reservation) error {
	return uc.ReservationRepository.UpdateReservationByID(reservation)
}

func (uc *ReservationUsecase) DeleteReservationByID(id int) error {
	return uc.ReservationRepository.DeleteReservationByID(id)
}

func (uc *ReservationUsecase) SearchPendingReservationsByVehicleID(vehicleID int) (*entities.Reservations, error) {
	return uc.ReservationRepository.SearchPendingReservationsByVehicleID(vehicleID)
}


func (uc *ReservationUsecase) freeParkingFromExpiredReservation(parkingID, reservationID int) error {
	parking, err := uc.ParkingRepository.SearchParkingByID(parkingID)
	if err != nil {
		return fmt.Errorf("error al obtener parking: %w", err)
	}

	
	
	if parking.IsActive {
		
		
		

		parking.IsActive = false
		err = uc.ParkingRepository.UpdateParkingByID(parking)
		if err != nil {
			return fmt.Errorf("error al liberar parking: %w", err)
		}
		log.Printf("Parking %d liberado por cancelación de reserva %d", parkingID, reservationID)
	}

	return nil
}
func (uc *ReservationUsecase) UpdateReservationStatus(id int, newStatus string) error {
	
	reservation, err := uc.SearchReservationByID(id)
	if err != nil {
		return fmt.Errorf("error al buscar reserva: %w", err)
	}

	if reservation == nil {
		return fmt.Errorf("reserva con ID %d no encontrada", id)
	}

	
	switch newStatus {
	case "cancelled":
		
		if reservation.Status != "pending" && reservation.Status != "active" {
			return fmt.Errorf("solo se pueden cancelar reservas pendientes o activas")
		}
	case "active":
		
		if reservation.Status != "pending" {
			return fmt.Errorf("solo se pueden activar reservas pendientes")
		}
	case "completed":
		
		if reservation.Status != "active" {
			return fmt.Errorf("solo se pueden completar reservas activas")
		}
	default:
		
	}

	
	reservation.Status = newStatus

	
	return uc.UpdateReservationById(reservation)
}




func (uc *ReservationUsecase) CancelExpiredPendingReservations() (int, error) {
	
	expiredReservations, err := uc.ReservationRepository.SearchReservationsByStatus([]string{"pending", "active"})
	if err != nil {
		return 0, fmt.Errorf("error al buscar reservas: %w", err)
	}

	if expiredReservations == nil || len(*expiredReservations) == 0 {
		return 0, nil
	}

	count := 0
	gracePeriod := 15 * time.Minute
	now := time.Now()

	for _, reservation := range *expiredReservations {
		
		if now.After(reservation.StartTime.Add(gracePeriod)) {
			
			reservation.Status = "cancelled"
			if err := uc.UpdateReservationById(&reservation); err != nil {
				log.Printf("Error al cancelar reserva expirada ID %d: %v", reservation.ID, err)
				continue
			}

			
			err = uc.freeParkingFromExpiredReservation(reservation.ParkingID, reservation.ID)
			if err != nil {
				log.Printf("Error al liberar parking de reserva expirada %d: %v", reservation.ID, err)
			}

			count++
			log.Printf("Reserva ID %d cancelada automáticamente por expiración", reservation.ID)
		}
	}

	return count, nil
}


func (uc *ReservationUsecase) SearchReservationsStartingWithin(timeLimit time.Time) (*entities.Reservations, error) {
	
	return uc.ReservationRepository.SearchReservationsStartingWithin(timeLimit)
}


func (uc *ReservationUsecase) ActivateReservationParking(reservationID int) error {
	
	reservation, err := uc.ReservationRepository.SearchReservationByID(reservationID)
	if err != nil {
		return fmt.Errorf("error al buscar reserva: %w", err)
	}

	if reservation == nil {
		return fmt.Errorf("reserva no encontrada")
	}

	
	parking, err := uc.ParkingRepository.SearchParkingByID(reservation.ParkingID)
	if err != nil {
		return fmt.Errorf("error al obtener parking: %w", err)
	}

	parking.IsActive = true 
	err = uc.ParkingRepository.UpdateParkingByID(parking)
	if err != nil {
		return fmt.Errorf("error al activar parking: %w", err)
	}

	
	reservation.Status = "active"
	err = uc.ReservationRepository.UpdateReservationByID(reservation)
	if err != nil {
		return fmt.Errorf("error al actualizar estado de reserva: %w", err)
	}

	return nil
}



func (uc *ReservationUsecase) findAvailableParkingForReservation(startTime, endTime time.Time, isImmediateReservation bool) (*parkingEntities.Parking, error) {
	
	allParkings, err := uc.ParkingRepository.SearchParkings()
	if err != nil {
		return nil, fmt.Errorf("error al obtener parkings: %w", err)
	}

	log.Printf("Buscando parking disponible para período %v a %v (inmediata: %v)", startTime, endTime, isImmediateReservation)
	log.Printf("Total de parkings a evaluar: %d", len(*allParkings))

	var unavailableReasons []string

	
	for _, parking := range *allParkings {
		log.Printf("Evaluando parking ID %d", parking.ID)

		available, err := uc.isParkingAvailableForReservation(parking.ID, startTime, endTime, isImmediateReservation)
		if err != nil {
			log.Printf("Error al evaluar parking %d: %v", parking.ID, err)
			unavailableReasons = append(unavailableReasons, fmt.Sprintf("Parking %d: error - %v", parking.ID, err))
			continue
		}

		if available {
			log.Printf("✅ Parking %d SELECCIONADO - está disponible", parking.ID)
			return &parking, nil
		} else {
			log.Printf("❌ Parking %d no disponible", parking.ID)
			unavailableReasons = append(unavailableReasons, fmt.Sprintf("Parking %d: no disponible", parking.ID))
		}
	}

	
	log.Printf("❌ NO hay parkings disponibles. Razones: %v", unavailableReasons)
	return nil, fmt.Errorf("no hay parkings disponibles para el período solicitado")
}


func (uc *ReservationUsecase) isParkingAvailableForReservation(
	parkingID int,
	startTime, endTime time.Time,
	isImmediate bool,
) (bool, error) {
	
	overlaps, err := uc.ReservationRepository.SearchOverlappingReservations(
		parkingID, startTime, endTime,
	)
	if err != nil {
		return false, err
	}
	if len(*overlaps) > 0 {
		return false, nil
	}

	
	parking, err := uc.ParkingRepository.SearchParkingByID(parkingID)
	if err != nil {
		return false, err
	}
	if parking.IsActive {
		return false, nil
	}

	return true, nil
}

func (uc *ReservationUsecase) freeParkingFromCancelledReservation(parkingID, reservationID int) error {
	
	

	hasActiveUsage, err := uc.parkingHasActiveUsage(parkingID)
	if err != nil {
		return err
	}

	if !hasActiveUsage {
		
		parking, err := uc.ParkingRepository.SearchParkingByID(parkingID)
		if err != nil {
			return err
		}

		parking.IsActive = false
		return uc.ParkingRepository.UpdateParkingByID(parking)
	}

	return nil
}

func (uc *ReservationUsecase) parkingHasActiveUsage(parkingID int) (bool, error) {
	
	
	return false, nil
}

package application

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gonzalohonorato/servercorego/core/parking/domain/repositories"
	"github.com/gonzalohonorato/servercorego/core/reservation/domain/entities"
)


type ReservationScheduler struct {
	usecase           *ReservationUsecase
	parkingRepository repositories.ParkingRepository
	stop              chan bool
	running           bool
}


func NewReservationScheduler(usecase *ReservationUsecase, parkingRepo repositories.ParkingRepository) *ReservationScheduler {
	return &ReservationScheduler{
		usecase:           usecase,
		parkingRepository: parkingRepo,
		stop:              make(chan bool),
		running:           false,
	}
}


func (s *ReservationScheduler) Start() {
	
	if s.running {
		log.Println("El scheduler de reservas ya está en ejecución")
		return
	}

	s.running = true
	log.Println("Iniciando scheduler completo de reservas...")

	
	s.runCancellationTask()
	s.runParkingActivationTask()

	
	go s.startCancellationScheduler()

	
	go s.startParkingActivationScheduler()
}


func (s *ReservationScheduler) startCancellationScheduler() {
	
	intervalStr := os.Getenv("RESERVATION_CANCELLATION_INTERVAL")
	interval := 10 * time.Minute 

	if intervalStr != "" {
		if minutes, err := strconv.Atoi(intervalStr); err == nil && minutes > 0 {
			interval = time.Duration(minutes) * time.Minute
		}
	}

	log.Printf("Iniciando scheduler de cancelación de reservas con intervalo de %v", interval)
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ticker.C:
			s.runCancellationTask()
		case <-s.stop:
			ticker.Stop()
			log.Println("Scheduler de cancelación de reservas detenido")
			return
		}
	}
}


func (s *ReservationScheduler) startParkingActivationScheduler() {
	
	intervalStr := os.Getenv("PARKING_ACTIVATION_INTERVAL")
	interval := 5 * time.Minute

	if intervalStr != "" {
		if minutes, err := strconv.Atoi(intervalStr); err == nil && minutes > 0 {
			interval = time.Duration(minutes) * time.Minute
		}
	}

	log.Printf("Iniciando scheduler de activación de parkings con intervalo de %v", interval)
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ticker.C:
			s.runParkingActivationTask()
		case <-s.stop:
			ticker.Stop()
			log.Println("Scheduler de activación de parkings detenido")
			return
		}
	}
}


func (s *ReservationScheduler) Stop() {
	if s.running {
		close(s.stop) 
		s.running = false
		log.Println("Deteniendo todos los schedulers de reservas...")
	}
}


func (s *ReservationScheduler) runCancellationTask() {
	log.Println("Ejecutando tarea de cancelación de reservas expiradas...")
	count, err := s.usecase.CancelExpiredPendingReservations()
	if err != nil {
		log.Printf("Error al ejecutar cancelación de reservas: %v", err)
	} else {
		if count > 0 {
			log.Printf("Proceso completado: %d reservas canceladas por expiración", count)
		}
	}
}


func (s *ReservationScheduler) runParkingActivationTask() {
	log.Println("Ejecutando tarea de activación de parkings para reservas próximas...")
	count, err := s.activateParkingsForUpcomingReservations()
	if err != nil {
		log.Printf("Error al activar parkings para reservas: %v", err)
	} else {
		if count > 0 {
			log.Printf("Proceso completado: %d parkings activados para reservas próximas", count)
		}
	}
}



func (s *ReservationScheduler) activateParkingsForUpcomingReservations() (int, error) {
	now := time.Now()
	oneHourFromNow := now.Add(1 * time.Hour)

	
	upcomingReservations, err := s.usecase.ReservationRepository.SearchReservationsStartingWithin(oneHourFromNow)
	if err != nil {
		return 0, err
	}

	if upcomingReservations == nil || len(*upcomingReservations) == 0 {
		return 0, nil
	}

	count := 0
	for _, reservation := range *upcomingReservations {
		
		if reservation.Status != "pending" {
			continue
		}

		
		timeUntilStart := time.Until(reservation.StartTime)
		if timeUntilStart <= time.Hour && timeUntilStart > 0 {
			err := s.activateReservationParking(reservation)
			if err != nil {
				log.Printf("Error al activar parking para reserva %d: %v", reservation.ID, err)
				continue
			}
			count++
			log.Printf("Parking activado para reserva ID %d, parking ID %d", reservation.ID, reservation.ParkingID)
		}
	}

	return count, nil
}


func (s *ReservationScheduler) activateReservationParking(reservation entities.Reservation) error {
	
	parking, err := s.parkingRepository.SearchParkingByID(reservation.ParkingID)
	if err != nil {
		return err
	}

	
	if parking.IsActive {
		log.Printf("Advertencia: Parking %d ya está ocupado, reserva %d puede tener conflicto",
			parking.ID, reservation.ID)
		return nil 
	}

	
	parking.IsActive = true
	err = s.parkingRepository.UpdateParkingByID(parking)
	if err != nil {
		return err
	}

	
	reservation.Status = "active"
	return s.usecase.UpdateReservationById(&reservation)
}

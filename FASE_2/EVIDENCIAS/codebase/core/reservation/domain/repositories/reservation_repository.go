package repositories

import (
	"time"

	"github.com/gonzalohonorato/servercorego/core/reservation/domain/entities"
)

type ReservationRepository interface {
	SearchReservationByID(id int) (*entities.Reservation, error)
	SearchReservations() (*entities.Reservations, error)
	SearchReservationsByDateAndStatus(date string, status []string) (*entities.Reservations, error)
	SearchReservationByUserIDAndDates(userID string, startDate string, endDate string) (*entities.Reservations, error)
	SearchPendingReservationsByVehicleID(vehicleID int) (*entities.Reservations, error)
	CreateReservation(reservation *entities.Reservation) error
	UpdateReservationByID(reservation *entities.Reservation) error
	DeleteReservationByID(id int) error
	SearchReservationsByStatus(statuses []string) (*entities.Reservations, error)
	SearchReservationsStartingWithin(timeLimit time.Time) (*entities.Reservations, error)

	SearchPendingReservationsByParkingAndTime(parkingID int, timeLimit time.Time) (*entities.Reservations, error)
	SearchActiveReservationByParkingAndTime(parkingID int, currentTime time.Time) (*entities.Reservation, error)
	SearchOverlappingReservations(parkingID int, startTime, endTime time.Time) (*entities.Reservations, error)
}

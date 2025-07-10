package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/gonzalohonorato/servercorego/core/reservation/domain/entities"
	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TimescaleReservationRepository struct {
	dbPool *pgxpool.Pool
}

func NewTimescaleDBRepository(pool *pgxpool.Pool) *TimescaleReservationRepository {
	return &TimescaleReservationRepository{
		dbPool: pool,
	}
}

func (r *TimescaleReservationRepository) SearchReservationByID(id int) (*entities.Reservation, error) {
	ctx := context.Background()
	query := `SELECT * FROM reservation WHERE id = $1`
	row := r.dbPool.QueryRow(ctx, query, id)
	var res entities.Reservation
	err := row.Scan(&res.ID, &res.CustomerID, &res.ParkingID, &res.VehicleID, &res.StartTime, &res.EndTime, &res.Status, &res.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *TimescaleReservationRepository) SearchReservations() (*entities.Reservations, error) {
	ctx := context.Background()
	query := `SELECT * FROM reservation`
	rows, err := r.dbPool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations entities.Reservations
	for rows.Next() {
		var res entities.Reservation
		if err := rows.Scan(&res.ID, &res.CustomerID, &res.ParkingID, &res.VehicleID, &res.StartTime, &res.EndTime, &res.Status, &res.CreatedAt); err != nil {
			return nil, err
		}
		reservations = append(reservations, res)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &reservations, nil
}
func (r *TimescaleReservationRepository) SearchReservationsByDateAndStatus(dateStr string, statuses []string) (*entities.Reservations, error) {
	ctx := context.Background()

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}

	
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	endOfDay := startOfDay.Add(24 * time.Hour)

	
	query := `
		SELECT * FROM reservation 
		WHERE start_time >= $1 AND start_time < $2 
		AND status = ANY($3)
	`

	rows, err := r.dbPool.Query(ctx, query, startOfDay, endOfDay, statuses)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations entities.Reservations
	for rows.Next() {
		var res entities.Reservation
		if err := rows.Scan(&res.ID, &res.CustomerID, &res.ParkingID, &res.VehicleID, &res.StartTime, &res.EndTime, &res.Status, &res.CreatedAt); err != nil {
			return nil, err
		}
		reservations = append(reservations, res)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &reservations, nil
}

func (r *TimescaleReservationRepository) CreateReservation(reservation *entities.Reservation) error {
	ctx := context.Background()
	query := `
	INSERT INTO reservation (
		customer_id, parking_id, vehicle_id, start_time, end_time, status, created_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7
	);
`
	_, err := r.dbPool.Exec(ctx, query,
		reservation.CustomerID,
		reservation.ParkingID,
		reservation.VehicleID,
		reservation.StartTime,
		reservation.EndTime,
		reservation.Status,
		time.Now())
	return err
}

func (r *TimescaleReservationRepository) UpdateReservationByID(res *entities.Reservation) error {
	ctx := context.Background()
	query := `UPDATE reservation SET customer_id = $2, parking_id = $3, vehicle_id = $4, start_time = $5, end_time = $6, status = $7 WHERE id = $1`

	_, err := r.dbPool.Exec(ctx, query,
		res.ID,
		res.CustomerID,
		res.ParkingID,
		res.VehicleID,
		res.StartTime,
		res.EndTime,
		res.Status)
	return err
}
func (r *TimescaleReservationRepository) SearchReservationByUserIDAndDates(userID string, startDate, endDate string) (*entities.Reservations, error) {
	ctx := context.Background()
	startTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, err
	}

	endTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, err
	}

	
	endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 999999999, time.Local)

	query := `
		SELECT * FROM reservation 
		WHERE customer_id = $1 
		AND created_at >= $2 
		AND created_at <= $3
	`
	fmt.Println("query", startTime)
	rows, err := r.dbPool.Query(ctx, query, userID, startTime, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations entities.Reservations
	for rows.Next() {
		var res entities.Reservation
		if err := rows.Scan(&res.ID, &res.CustomerID, &res.ParkingID, &res.VehicleID, &res.StartTime, &res.EndTime, &res.Status, &res.CreatedAt); err != nil {
			return nil, err
		}
		reservations = append(reservations, res)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &reservations, nil
}
func (r *TimescaleReservationRepository) DeleteReservationByID(id int) error {
	ctx := context.Background()
	query := `DELETE FROM reservation WHERE id = $1`
	_, err := r.dbPool.Exec(ctx, query, id)
	return err
}

func (r *TimescaleReservationRepository) SearchPendingReservationsByVehicleID(vehicleID int) (*entities.Reservations, error) {
	ctx := context.Background()
	now := time.Now()

	query := `
        SELECT * FROM reservation 
        WHERE vehicle_id = $1 
        AND status = 'pending'
        AND start_time <= $2 
        AND end_time >= $2
        ORDER BY start_time ASC
    `

	rows, err := r.dbPool.Query(ctx, query, vehicleID, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations entities.Reservations
	for rows.Next() {
		var res entities.Reservation
		if err := rows.Scan(&res.ID, &res.CustomerID, &res.ParkingID, &res.VehicleID, &res.StartTime, &res.EndTime, &res.Status, &res.CreatedAt); err != nil {
			return nil, err
		}
		reservations = append(reservations, res)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &reservations, nil
}


func (r *TimescaleReservationRepository) SearchReservationsByStatus(statuses []string) (*entities.Reservations, error) {
	ctx := context.Background()

	query := `
        SELECT id, customer_id, parking_id, vehicle_id, start_time, end_time, status, created_at
        FROM reservation 
        WHERE status = ANY($1)
    `

	rows, err := r.dbPool.Query(ctx, query, statuses)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations entities.Reservations
	for rows.Next() {
		var res entities.Reservation
		if err := rows.Scan(&res.ID, &res.CustomerID, &res.ParkingID, &res.VehicleID, &res.StartTime, &res.EndTime, &res.Status, &res.CreatedAt); err != nil {
			return nil, err
		}
		reservations = append(reservations, res)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &reservations, nil
}


func (r *TimescaleReservationRepository) SearchReservationsStartingWithin(timeLimit time.Time) (*entities.Reservations, error) {
	ctx := context.Background()

	query := `
        SELECT id, customer_id, parking_id, vehicle_id, start_time, end_time, status, created_at 
        FROM reservation 
        WHERE status = 'pending' 
        AND start_time <= $1 
        AND start_time > NOW()
        ORDER BY start_time ASC
    `

	rows, err := r.dbPool.Query(ctx, query, timeLimit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations entities.Reservations
	for rows.Next() {
		var res entities.Reservation
		if err := rows.Scan(&res.ID, &res.CustomerID, &res.ParkingID, &res.VehicleID, &res.StartTime, &res.EndTime, &res.Status, &res.CreatedAt); err != nil {
			return nil, err
		}
		reservations = append(reservations, res)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &reservations, nil
}


func (r *TimescaleReservationRepository) SearchOverlappingReservations(
	parkingID int,
	startTime, endTime time.Time,
) (*entities.Reservations, error) {
	ctx := context.Background()
	query := `
      SELECT id, customer_id, parking_id, vehicle_id, start_time, end_time, status, created_at
      FROM reservation
      WHERE parking_id = $1
        AND status IN ('pending','active')
        AND NOT (end_time <= $2 OR start_time >= $3)
    `
	rows, err := r.dbPool.Query(ctx, query, parkingID, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list entities.Reservations
	for rows.Next() {
		var r entities.Reservation
		if err := rows.Scan(
			&r.ID, &r.CustomerID, &r.ParkingID, &r.VehicleID,
			&r.StartTime, &r.EndTime, &r.Status, &r.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, r)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return &list, nil
}


func (r *TimescaleReservationRepository) SearchPendingReservationsByParkingAndTime(parkingID int, timeLimit time.Time) (*entities.Reservations, error) {
	ctx := context.Background()

	query := `
        SELECT id, customer_id, parking_id, vehicle_id, start_time, end_time, status, created_at 
        FROM reservation 
        WHERE parking_id = $1 
        AND status = 'pending'
        AND start_time <= $2 
        AND start_time > NOW()
        ORDER BY start_time ASC
    `

	rows, err := r.dbPool.Query(ctx, query, parkingID, timeLimit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations entities.Reservations
	for rows.Next() {
		var res entities.Reservation
		if err := rows.Scan(&res.ID, &res.CustomerID, &res.ParkingID, &res.VehicleID, &res.StartTime, &res.EndTime, &res.Status, &res.CreatedAt); err != nil {
			return nil, err
		}
		reservations = append(reservations, res)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &reservations, nil
}


func (r *TimescaleReservationRepository) SearchActiveReservationByParkingAndTime(parkingID int, currentTime time.Time) (*entities.Reservation, error) {
	ctx := context.Background()

	query := `
        SELECT id, customer_id, parking_id, vehicle_id, start_time, end_time, status, created_at 
        FROM reservation 
        WHERE parking_id = $1 
        AND status IN ('active')
        AND start_time <= $2 
        AND end_time >= $2
        LIMIT 1
    `

	var res entities.Reservation
	err := r.dbPool.QueryRow(ctx, query, parkingID, currentTime).Scan(
		&res.ID, &res.CustomerID, &res.ParkingID, &res.VehicleID,
		&res.StartTime, &res.EndTime, &res.Status, &res.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil 
		}
		return nil, err
	}

	return &res, nil
}

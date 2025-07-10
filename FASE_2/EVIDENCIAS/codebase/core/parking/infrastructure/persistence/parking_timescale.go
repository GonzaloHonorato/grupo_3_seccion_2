package persistence

import (
	"context"

	"github.com/gonzalohonorato/servercorego/core/parking/domain/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TimescaleParkingRepository struct {
	dbPool *pgxpool.Pool
}

func NewTimescaleDBRepository(pool *pgxpool.Pool) *TimescaleParkingRepository {
	return &TimescaleParkingRepository{
		dbPool: pool,
	}
}

func (r *TimescaleParkingRepository) SearchParkingByID(id int) (*entities.Parking, error) {
	ctx := context.Background()
	query := `SELECT * FROM parking WHERE id = $1`
	row := r.dbPool.QueryRow(ctx, query, id)
	var p entities.Parking
	err := row.Scan(&p.ID, &p.Code, &p.Location, &p.Zone, &p.IsActive)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *TimescaleParkingRepository) SearchParkings() (*entities.Parkings, error) {
	ctx := context.Background()
	query := `SELECT * FROM parking`
	rows, err := r.dbPool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parkings entities.Parkings
	for rows.Next() {
		var p entities.Parking
		if err := rows.Scan(&p.ID, &p.Code, &p.Location, &p.Zone, &p.IsActive); err != nil {
			return nil, err
		}
		parkings = append(parkings, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &parkings, nil
}
func (r *TimescaleParkingRepository) SearchAvailableParkings() (*entities.Parkings, error) {
	ctx := context.Background()

	
	query := `
        SELECT p.id, p.code, p.location, p.zone, p.is_active 
        FROM parking p
        WHERE p.is_active = false
    `

	rows, err := r.dbPool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parkings entities.Parkings
	for rows.Next() {
		var p entities.Parking
		if err := rows.Scan(&p.ID, &p.Code, &p.Location, &p.Zone, &p.IsActive); err != nil {
			return nil, err
		}
		parkings = append(parkings, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &parkings, nil
}
func (r *TimescaleParkingRepository) CreateParking(parking *entities.Parking) error {
	ctx := context.Background()
	query := `
	INSERT INTO parkings (
		code, location, zone, is_active
	) VALUES (
		$1, $2, $3, $4
	);
`
	_, err := r.dbPool.Exec(ctx, query, parking.Code, parking.Location, parking.Zone, parking.IsActive)
	return err
}

func (r *TimescaleParkingRepository) UpdateParkingByID(p *entities.Parking) error {
	ctx := context.Background()
	query := `UPDATE parking SET code = $2, location = $3, zone = $4, is_active = $5 WHERE id = $1`

	_, err := r.dbPool.Exec(ctx, query, p.ID, p.Code, p.Location, p.Zone, p.IsActive)
	return err
}

func (r *TimescaleParkingRepository) DeleteParkingByID(id int) error {
	ctx := context.Background()
	query := `DELETE FROM parking WHERE id = $1`
	_, err := r.dbPool.Exec(ctx, query, id)
	return err
}

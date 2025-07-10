package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/gonzalohonorato/servercorego/core/vehicle/domain/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TimescaleVehicleRepository struct {
	dbPool *pgxpool.Pool
}

func NewTimescaleVehicleRepository(pool *pgxpool.Pool) *TimescaleVehicleRepository {
	return &TimescaleVehicleRepository{
		dbPool: pool,
	}
}

func (r *TimescaleVehicleRepository) SearchVehicleByID(id int) (*entities.Vehicle, error) {
	ctx := context.Background()
	fmt.Println("Searching vehicle by ID:", id)
	query := `SELECT id, plate, brand, model, vehicle_type, customer_id, created_at FROM vehicle WHERE id = $1`
	row := r.dbPool.QueryRow(ctx, query, id)
	var v entities.Vehicle
	err := row.Scan(&v.ID, &v.Plate, &v.Brand, &v.Model, &v.VehicleType, &v.CustomerID, &v.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &v, nil
}
func (r *TimescaleVehicleRepository) SearchVehiclesByCustomerID(customerID string) (*entities.Vehicles, error) {
	ctx := context.Background()
	query := `SELECT id, plate, brand, model, vehicle_type, customer_id, created_at FROM vehicle WHERE customer_id = $1`
	rows, err := r.dbPool.Query(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var vehicles entities.Vehicles
	for rows.Next() {
		var v entities.Vehicle
		if err := rows.Scan(&v.ID, &v.Plate, &v.Brand, &v.Model, &v.VehicleType, &v.CustomerID, &v.CreatedAt); err != nil {
			return nil, err
		}
		vehicles = append(vehicles, v)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &vehicles, nil
}

func (r *TimescaleVehicleRepository) SearchVehicleByPlate(plate string) (*entities.Vehicle, error) {
	ctx := context.Background()
	query := `SELECT id, plate, brand, model, vehicle_type, customer_id, created_at FROM vehicle WHERE plate = $1`
	row := r.dbPool.QueryRow(ctx, query, plate)
	var v entities.Vehicle
	err := row.Scan(&v.ID, &v.Plate, &v.Brand, &v.Model, &v.VehicleType, &v.CustomerID, &v.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *TimescaleVehicleRepository) SearchVehicles() (*entities.Vehicles, error) {
	ctx := context.Background()
	query := `SELECT id, plate, brand, model, vehicle_type, customer_id, created_at FROM vehicle`
	rows, err := r.dbPool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vehicles entities.Vehicles
	for rows.Next() {
		var v entities.Vehicle
		if err := rows.Scan(&v.ID, &v.Plate, &v.Brand, &v.Model, &v.VehicleType, &v.CustomerID, &v.CreatedAt); err != nil {
			return nil, err
		}
		vehicles = append(vehicles, v)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &vehicles, nil
}

func (r *TimescaleVehicleRepository) CreateVehicle(vehicle *entities.Vehicle) error {
	ctx := context.Background()
	query := `
	INSERT INTO vehicle (
		plate, brand, model, vehicle_type, customer_id, created_at
	) VALUES (
		$1, $2, $3, $4, $5, $6
	);
`
	
	if vehicle.CreatedAt.IsZero() {
		vehicle.CreatedAt = time.Now()
	}

	_, err := r.dbPool.Exec(ctx, query, vehicle.Plate, vehicle.Brand, vehicle.Model, vehicle.VehicleType, vehicle.CustomerID, vehicle.CreatedAt)
	return err
}

func (r *TimescaleVehicleRepository) UpdateVehicleByID(v *entities.Vehicle) error {
	ctx := context.Background()
	query := `UPDATE vehicle SET plate = $2, brand = $3, model = $4, vehicle_type = $5, customer_id = $6, created_at = $7 WHERE id = $1`

	_, err := r.dbPool.Exec(ctx, query, v.ID, v.Plate, v.Brand, v.Model, v.VehicleType, v.CustomerID, v.CreatedAt)
	return err
}

func (r *TimescaleVehicleRepository) DeleteVehicleByID(id int) error {
	ctx := context.Background()
	query := `DELETE FROM vehicle WHERE id = $1`
	_, err := r.dbPool.Exec(ctx, query, id)
	return err
}

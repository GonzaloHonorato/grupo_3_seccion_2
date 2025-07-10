package persistence

import (
	"context"
	"fmt"
	"strings"

	"github.com/gonzalohonorato/servercorego/core/parkingusage/domain/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TimescaleParkingUsageRepository struct {
	dbPool *pgxpool.Pool
}

func NewTimescaleDBRepository(pool *pgxpool.Pool) *TimescaleParkingUsageRepository {
	return &TimescaleParkingUsageRepository{
		dbPool: pool,
	}
}

func (r *TimescaleParkingUsageRepository) SearchParkingUsageByID(id int) (*entities.ParkingUsage, error) {
	ctx := context.Background()
	query := `SELECT id, reservation_id, vehicle_id, parking_id, entry_time, exit_time, 
			  ocr_plate, qr_scanned, registered_by, manual_entry, visitor_name, visitor_rut, 
			  visitor_contact, zone FROM parking_usage WHERE id = $1`
	row := r.dbPool.QueryRow(ctx, query, id)
	var p entities.ParkingUsage
	err := row.Scan(
		&p.ID, &p.ReservationID, &p.VehicleID, &p.ParkingID, &p.EntryTime, &p.ExitTime,
		&p.OcrPlate, &p.QrScanned, &p.RegisteredBy, &p.ManualEntry, &p.VisitorName,
		&p.VisitorRut, &p.VisitorContact, &p.Zone,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *TimescaleParkingUsageRepository) SearchParkingUsages() (*entities.ParkingUsages, error) {
	ctx := context.Background()
	query := `SELECT id, reservation_id, vehicle_id, parking_id, entry_time, exit_time, 
              ocr_plate, qr_scanned, registered_by, manual_entry, visitor_name, visitor_rut, 
              visitor_contact, zone FROM parking_usage`
	rows, err := r.dbPool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	parkingUsages := entities.ParkingUsages{} 
	for rows.Next() {
		var p entities.ParkingUsage
		if err := rows.Scan(
			&p.ID, &p.ReservationID, &p.VehicleID, &p.ParkingID, &p.EntryTime, &p.ExitTime,
			&p.OcrPlate, &p.QrScanned, &p.RegisteredBy, &p.ManualEntry, &p.VisitorName,
			&p.VisitorRut, &p.VisitorContact, &p.Zone,
		); err != nil {
			return nil, err
		}
		parkingUsages = append(parkingUsages, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &parkingUsages, nil 
}
func (r *TimescaleParkingUsageRepository) SearchActiveParkingUsages() (*entities.ParkingUsages, error) {
	ctx := context.Background()
	query := `SELECT id, reservation_id, vehicle_id, parking_id, entry_time, exit_time, 
              ocr_plate, qr_scanned, registered_by, manual_entry, visitor_name, visitor_rut, 
              visitor_contact, zone FROM active_parking_usages`
	rows, err := r.dbPool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	parkingUsages := entities.ParkingUsages{} 
	for rows.Next() {
		var p entities.ParkingUsage
		if err := rows.Scan(
			&p.ID, &p.ReservationID, &p.VehicleID, &p.ParkingID, &p.EntryTime, &p.ExitTime,
			&p.OcrPlate, &p.QrScanned, &p.RegisteredBy, &p.ManualEntry, &p.VisitorName,
			&p.VisitorRut, &p.VisitorContact, &p.Zone,
		); err != nil {
			return nil, err
		}
		parkingUsages = append(parkingUsages, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &parkingUsages, nil 
}
func (r *TimescaleParkingUsageRepository) SearchParkingUsagesByVehicleID(id int) (*entities.ParkingUsages, error) {
	ctx := context.Background()
	query := `SELECT id, reservation_id, vehicle_id, parking_id, entry_time, exit_time,
			  ocr_plate, qr_scanned, registered_by, manual_entry, visitor_name, visitor_rut,	
			  visitor_contact, zone FROM parking_usage WHERE vehicle_id = $1`
	rows, err := r.dbPool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	parkingUsages := entities.ParkingUsages{} 
	for rows.Next() {
		var p entities.ParkingUsage
		if err := rows.Scan(
			&p.ID, &p.ReservationID, &p.VehicleID, &p.ParkingID, &p.EntryTime, &p.ExitTime,
			&p.OcrPlate, &p.QrScanned, &p.RegisteredBy, &p.ManualEntry, &p.VisitorName,
			&p.VisitorRut, &p.VisitorContact, &p.Zone,
		); err != nil {
			return nil, err
		}
		parkingUsages = append(parkingUsages, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &parkingUsages, nil 
}
func (r *TimescaleParkingUsageRepository) CreateParkingUsage(parkingUsage *entities.ParkingUsage) error {
	ctx := context.Background()
	query := `
	INSERT INTO parking_usage (
		reservation_id, vehicle_id, parking_id, entry_time, exit_time, 
		ocr_plate, qr_scanned, registered_by, manual_entry, visitor_name, 
		visitor_rut, visitor_contact, zone
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
	) RETURNING id;
`
	err := r.dbPool.QueryRow(ctx, query,
		parkingUsage.ReservationID, parkingUsage.VehicleID, parkingUsage.ParkingID,
		parkingUsage.EntryTime, parkingUsage.ExitTime, parkingUsage.OcrPlate,
		parkingUsage.QrScanned, parkingUsage.RegisteredBy, parkingUsage.ManualEntry,
		parkingUsage.VisitorName, parkingUsage.VisitorRut, parkingUsage.VisitorContact,
		parkingUsage.Zone).Scan(&parkingUsage.ID)
	return err
}

func (r *TimescaleParkingUsageRepository) UpdateParkingUsageByID(p *entities.ParkingUsage) error {
	ctx := context.Background()
	query := `UPDATE parking_usage SET 
		reservation_id = $2, vehicle_id = $3, parking_id = $4, entry_time = $5, 
		exit_time = $6, ocr_plate = $7, qr_scanned = $8, registered_by = $9, 
		manual_entry = $10, visitor_name = $11, visitor_rut = $12, 
		visitor_contact = $13, zone = $14 
	WHERE id = $1`

	_, err := r.dbPool.Exec(ctx, query,
		p.ID, p.ReservationID, p.VehicleID, p.ParkingID, p.EntryTime, p.ExitTime,
		p.OcrPlate, p.QrScanned, p.RegisteredBy, p.ManualEntry, p.VisitorName,
		p.VisitorRut, p.VisitorContact, p.Zone)
	return err
}

func (r *TimescaleParkingUsageRepository) DeleteParkingUsageByID(id int) error {
	ctx := context.Background()
	query := `DELETE FROM parking_usage WHERE id = $1`
	_, err := r.dbPool.Exec(ctx, query, id)
	return err
}



func (r *TimescaleParkingUsageRepository) SearchParkingUsagesByVehicleIDs(vehicleIDs []int, filters map[string]interface{}) (*[]entities.ParkingUsage, error) {
	ctx := context.Background()

	
	placeholders := make([]string, len(vehicleIDs))
	args := make([]interface{}, len(vehicleIDs))

	for i, id := range vehicleIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	
	baseQuery := fmt.Sprintf(`
        SELECT id, reservation_id, vehicle_id, parking_id, entry_time, exit_time,
        ocr_plate, qr_scanned, registered_by, manual_entry, visitor_name, visitor_rut,
        visitor_contact, zone 
        FROM parking_usage 
        WHERE vehicle_id IN (%s)
    `, strings.Join(placeholders, ","))

	
	var conditions []string
	paramCounter := len(vehicleIDs) + 1 

	
	if parkingID, ok := filters["parkingId"].(int); ok && parkingID > 0 {
		conditions = append(conditions, fmt.Sprintf("parking_id = $%d", paramCounter))
		args = append(args, parkingID)
		paramCounter++
	}

	if status, ok := filters["status"].(string); ok && status != "" {
		if status == "active" {
			conditions = append(conditions, "exit_time IS NULL")
		} else if status == "completed" {
			conditions = append(conditions, "exit_time IS NOT NULL")
		}
	}

	if startDate, ok := filters["startDate"].(string); ok && startDate != "" {
		conditions = append(conditions, fmt.Sprintf("entry_time >= $%d", paramCounter))
		args = append(args, startDate)
		paramCounter++
	}

	if endDate, ok := filters["endDate"].(string); ok && endDate != "" {
		conditions = append(conditions, fmt.Sprintf("entry_time <= $%d", paramCounter))
		args = append(args, endDate)
		paramCounter++
	}

	
	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	
	baseQuery += " ORDER BY entry_time DESC"

	
	rows, err := r.dbPool.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar consulta: %w", err)
	}
	defer rows.Close()

	
	parkingUsages := []entities.ParkingUsage{} 
	for rows.Next() {
		var p entities.ParkingUsage
		if err := rows.Scan(
			&p.ID, &p.ReservationID, &p.VehicleID, &p.ParkingID, &p.EntryTime, &p.ExitTime,
			&p.OcrPlate, &p.QrScanned, &p.RegisteredBy, &p.ManualEntry, &p.VisitorName,
			&p.VisitorRut, &p.VisitorContact, &p.Zone,
		); err != nil {
			return nil, fmt.Errorf("error al escanear fila: %w", err)
		}
		parkingUsages = append(parkingUsages, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error en rows: %w", err)
	}

	return &parkingUsages, nil 
}

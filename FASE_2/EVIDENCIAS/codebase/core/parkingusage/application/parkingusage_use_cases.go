package application

import (
	"fmt"
	"time"

	parkingEntity "github.com/gonzalohonorato/servercorego/core/parking/domain/entities"
	parkingRepository "github.com/gonzalohonorato/servercorego/core/parking/domain/repositories"
	"github.com/gonzalohonorato/servercorego/core/parkingusage/domain/entities"
	"github.com/gonzalohonorato/servercorego/core/parkingusage/domain/repositories"
	reservationEntity "github.com/gonzalohonorato/servercorego/core/reservation/domain/entities"
	reservationRepository "github.com/gonzalohonorato/servercorego/core/reservation/domain/repositories"
	vehicleEntity "github.com/gonzalohonorato/servercorego/core/vehicle/domain/entities"
	vehicleRepository "github.com/gonzalohonorato/servercorego/core/vehicle/domain/repositories"
	"github.com/gonzalohonorato/servercorego/core/websocket/infrastructure"
)

type ParkingUsageUsecase struct {
	ParkingUsageRepository repositories.ParkingUsageRepository
	ParkingRepository      parkingRepository.ParkingRepository
	VehicleRepository      vehicleRepository.VehicleRepository
	ReservationRepository  reservationRepository.ReservationRepository
	WebSocketService       *infrastructure.WebSocketService
}

func NewParkingUsageUsecase(
	parkingUsageRepo repositories.ParkingUsageRepository,
	parkingRepo parkingRepository.ParkingRepository,
	vehicleRepo vehicleRepository.VehicleRepository,
	reservationRepo reservationRepository.ReservationRepository,
	wsService *infrastructure.WebSocketService,
) *ParkingUsageUsecase {
	return &ParkingUsageUsecase{
		ParkingUsageRepository: parkingUsageRepo,
		ParkingRepository:      parkingRepo,
		VehicleRepository:      vehicleRepo,
		ReservationRepository:  reservationRepo,
		WebSocketService:       wsService,
	}
}

func (uc *ParkingUsageUsecase) SearchParkingUsageByID(id int) (*entities.ParkingUsage, error) {
	return uc.ParkingUsageRepository.SearchParkingUsageByID(id)
}

func (uc *ParkingUsageUsecase) SearchParkingUsages() (*entities.ParkingUsages, error) {
	return uc.ParkingUsageRepository.SearchParkingUsages()
}
func (uc *ParkingUsageUsecase) SearchParkingUsageByVehicleID(vehicleID int) (*entities.ParkingUsages, error) {
	return uc.ParkingUsageRepository.SearchParkingUsagesByVehicleID(vehicleID)
}
func (uc *ParkingUsageUsecase) SearchActiveParkingUsages() (*entities.ParkingUsages, error) {
	return uc.ParkingUsageRepository.SearchActiveParkingUsages()
}

func (uc *ParkingUsageUsecase) CreateParkingUsage(parkingUsage *entities.ParkingUsage) error {
	
	if parkingUsage.ParkingID == 0 {
		availableParking, err := uc.findAvailableParking()
		if err != nil {
			return fmt.Errorf("no hay estacionamientos disponibles: %w", err)
		}
		parkingUsage.ParkingID = availableParking.ID
	} else {
		
		isAvailable, reservationID, err := uc.isParkingAvailableForImmediate(parkingUsage.ParkingID)
		if err != nil {
			return fmt.Errorf("error al verificar disponibilidad: %w", err)
		}

		if !isAvailable {
			return fmt.Errorf("el estacionamiento no está disponible (conflicto con reservas próximas)")
		}

		
		if reservationID != nil {
			parkingUsage.ReservationID = reservationID
		}
	}

	
	if parkingUsage.EntryTime == nil {
		now := time.Now()
		parkingUsage.EntryTime = &now
	}

	
	err := uc.ParkingUsageRepository.CreateParkingUsage(parkingUsage)
	if err != nil {
		return err
	}

	
	if parkingUsage.ReservationID != nil {
		err = uc.consumeReservation(*parkingUsage.ReservationID)
		if err != nil {
			fmt.Printf("Advertencia: No se pudo consumir la reserva %d: %v\n", *parkingUsage.ReservationID, err)
		}
	}

	
	parking, err := uc.ParkingRepository.SearchParkingByID(parkingUsage.ParkingID)
	if err != nil {
		return fmt.Errorf("error al obtener el parking: %w", err)
	}

	if !parking.IsActive {
		parking.IsActive = true
		err = uc.ParkingRepository.UpdateParkingByID(parking)
		if err != nil {
			return fmt.Errorf("error al actualizar estado del parking: %w", err)
		}
	}

	if isRealtimeEntry(parkingUsage.EntryTime) {
		uc.notifyRealtimeEntry(parkingUsage)
	}

	return nil
}

type ExitRequest struct {
	ParkingUsageID int    `json:"parkingUsageId,omitempty"`
	Plate          string `json:"plate,omitempty"`
	ExitType       string `json:"exitType"` 
}

type ExitResponse struct {
	Success      bool                   `json:"success"`
	Message      string                 `json:"message"`
	ParkingUsage *entities.ParkingUsage `json:"parkingUsage,omitempty"`
	ErrorCode    string                 `json:"errorCode,omitempty"`
}


func (uc *ParkingUsageUsecase) ProcessParkingExit(request *ExitRequest) (*ExitResponse, error) {
	var parkingUsage *entities.ParkingUsage
	var err error

	
	switch request.ExitType {
	case "id":
		if request.ParkingUsageID == 0 {
			return &ExitResponse{
				Success:   false,
				Message:   "ID de parking usage requerido",
				ErrorCode: "PARKING_USAGE_ID_REQUIRED",
			}, nil
		}
		parkingUsage, err = uc.ParkingUsageRepository.SearchParkingUsageByID(request.ParkingUsageID)
	case "plate":
		fmt.Println("Procesando salida por patente:", request.Plate)
		if request.Plate == "" {
			return &ExitResponse{
				Success:   false,
				Message:   "Patente requerida",
				ErrorCode: "PLATE_REQUIRED",
			}, nil
		}
		parkingUsage, err = uc.findActiveParkingUsageByPlate(request.Plate)
	default:
		return &ExitResponse{
			Success:   false,
			Message:   "Tipo de salida no válido",
			ErrorCode: "INVALID_EXIT_TYPE",
		}, nil
	}

	if err != nil {
		return nil, fmt.Errorf("error al buscar parking usage: %w", err)
	}

	if parkingUsage == nil {
		response := &ExitResponse{
			Success:   false,
			Message:   "No se encontró un uso de estacionamiento activo",
			ErrorCode: "ACTIVE_PARKING_USAGE_NOT_FOUND",
		}
		uc.notifyExitRejection(response, request)
		return response, nil
	}

	
	if parkingUsage.ExitTime != nil {
		response := &ExitResponse{
			Success:   false,
			Message:   "Este registro ya tiene una hora de salida registrada",
			ErrorCode: "EXIT_ALREADY_REGISTERED",
		}
		uc.notifyExitRejection(response, request)
		return response, nil
	}

	
	now := time.Now()
	parkingUsage.ExitTime = &now

	
	err = uc.ParkingUsageRepository.UpdateParkingUsageByID(parkingUsage)
	if err != nil {
		return nil, fmt.Errorf("error al actualizar registro: %w", err)
	}

	
	parking, err := uc.ParkingRepository.SearchParkingByID(parkingUsage.ParkingID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener el parking: %w", err)
	}

	parking.IsActive = false
	err = uc.ParkingRepository.UpdateParkingByID(parking)
	if err != nil {
		return nil, fmt.Errorf("error al actualizar estado del parking: %w", err)
	}

	
	response := &ExitResponse{
		Success:      true,
		Message:      "Salida registrada exitosamente",
		ParkingUsage: parkingUsage,
	}

	uc.notifyExitSuccess(response, request)
	return response, nil
}


func (uc *ParkingUsageUsecase) findActiveParkingUsageByPlate(plate string) (*entities.ParkingUsage, error) {
	
	vehicle, err := uc.findVehicleByPlate(plate)
	if err != nil {
		return nil, fmt.Errorf("error al buscar vehículo por patente %s: %w", plate, err)
	}

	if vehicle == nil {
		
		return nil, nil
	}

	
	activeUsages, err := uc.ParkingUsageRepository.SearchActiveParkingUsages()
	if err != nil {
		return nil, fmt.Errorf("error al buscar registros activos: %w", err)
	}

	
	for _, usage := range *activeUsages {
		if usage.VehicleID != nil && *usage.VehicleID == vehicle.ID && usage.ExitTime == nil {
			return &usage, nil
		}
	}

	
	
	for _, usage := range *activeUsages {
		if usage.ExitTime == nil && usage.OcrPlate == plate {
			return &usage, nil
		}
	}

	return nil, nil
}
func (uc *ParkingUsageUsecase) isParkingAvailable(parkingID int) (bool, error) {
	
	parking, err := uc.ParkingRepository.SearchParkingByID(parkingID)
	if err != nil {
		return false, fmt.Errorf("error al obtener parking: %w", err)
	}

	
	if parking.IsActive {
		return false, nil
	}

	
	now := time.Now()
	sixHoursFromNow := now.Add(6 * time.Hour)

	hasPendingReservation, err := uc.parkingHasPendingReservationWithin(parkingID, sixHoursFromNow)
	if err != nil {
		return false, err
	}

	if hasPendingReservation {
		return false, nil 
	}

	
	hasActiveUsage, err := uc.parkingHasActiveUsage(parkingID)
	if err != nil {
		return false, err
	}

	if hasActiveUsage {
		return false, nil
	}

	return true, nil
}


func (uc *ParkingUsageUsecase) findAvailableParking() (*parkingEntity.Parking, error) {
	allParkings, err := uc.ParkingRepository.SearchParkings()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	sixHoursFromNow := now.Add(6 * time.Hour)

	for _, parking := range *allParkings {
		
		if parking.IsActive {
			continue 
		}

		
		hasPendingReservationSoon, err := uc.parkingHasPendingReservationWithin(parking.ID, sixHoursFromNow)
		if err != nil {
			fmt.Printf("Error al verificar reservas para parking %d: %v", parking.ID, err)
			continue
		}

		
		if hasPendingReservationSoon {
			continue
		}

		
		hasActiveUsage, err := uc.parkingHasActiveUsage(parking.ID)
		if err != nil {
			fmt.Printf("Error al verificar uso activo para parking %d: %v", parking.ID, err)
			continue
		}

		if hasActiveUsage {
			continue
		}

		
		return &parking, nil
	}

	return nil, fmt.Errorf("no hay parkings disponibles sin conflictos de reserva")
}
func (uc *ParkingUsageUsecase) isParkingAvailableForImmediate(parkingID int) (bool, *int, error) {
	
	parking, err := uc.ParkingRepository.SearchParkingByID(parkingID)
	if err != nil {
		return false, nil, err
	}

	if parking.IsActive {
		
		
		now := time.Now()
		reservationID, err := uc.getActiveReservationForParking(parkingID, now)
		if err != nil {
			return false, nil, err
		}

		if reservationID != nil {
			
			return true, reservationID, nil
		}

		
		return false, nil, nil
	}

	
	now := time.Now()
	sixHoursFromNow := now.Add(6 * time.Hour)
	hasPendingReservationSoon, err := uc.parkingHasPendingReservationWithin(parkingID, sixHoursFromNow)
	if err != nil {
		return false, nil, err
	}

	if hasPendingReservationSoon {
		return false, nil, nil 
	}

	
	return true, nil, nil
}


func (uc *ParkingUsageUsecase) parkingHasPendingReservationWithin(parkingID int, timeLimit time.Time) (bool, error) {
	reservations, err := uc.ReservationRepository.SearchPendingReservationsByParkingAndTime(parkingID, timeLimit)
	if err != nil {
		return false, err
	}

	if reservations == nil || len(*reservations) == 0 {
		return false, nil
	}

	return true, nil
}

func isRealtimeEntry(entryTime *time.Time) bool {
	if entryTime == nil {
		return false
	}

	now := time.Now()
	fiveMinAgo := now.Add(-5 * time.Minute)
	fiveMinAhead := now.Add(5 * time.Minute)

	return entryTime.After(fiveMinAgo) && entryTime.Before(fiveMinAhead)
}

func (uc *ParkingUsageUsecase) notifyRealtimeEntry(parkingUsage *entities.ParkingUsage) {
	if uc.WebSocketService != nil {
		uc.WebSocketService.BroadcastParkingUsage(parkingUsage)
	}
}

func (uc *ParkingUsageUsecase) UpdateParkingUsageById(parkingUsage *entities.ParkingUsage) error {
	return uc.ParkingUsageRepository.UpdateParkingUsageByID(parkingUsage)
}

func (uc *ParkingUsageUsecase) DeleteParkingUsageByID(id int) error {
	return uc.ParkingUsageRepository.DeleteParkingUsageByID(id)
}


type EntryRequest struct {
	
	ParkingID int    `json:"parkingId,omitempty"`
	EntryType string `json:"entryType"` 

	
	Plate string `json:"plate,omitempty"`

	
	UserID    string `json:"userId,omitempty"`
	VehicleID int    `json:"vehicleId,omitempty"`

	
	VisitorName    string `json:"visitorName,omitempty"`
	VisitorRut     string `json:"visitorRut,omitempty"`
	VisitorContact string `json:"visitorContact,omitempty"`
	Zone           string `json:"zone,omitempty"`
	RegisteredBy   string `json:"registeredBy,omitempty"`
}

type EntryResponse struct {
	Success      bool                   `json:"success"`
	Message      string                 `json:"message"`
	ParkingUsage *entities.ParkingUsage `json:"parkingUsage,omitempty"`
	ErrorCode    string                 `json:"errorCode,omitempty"`
}


func (uc *ParkingUsageUsecase) ProcessParkingEntry(request *EntryRequest) (*EntryResponse, error) {
	
	switch request.EntryType {
	case "ocr":
		return uc.processOCREntry(request)
	case "qr", "app":
		return uc.processQREntry(request)
	case "manual":
		return uc.processManualEntry(request)
	default:
		return &EntryResponse{
			Success:   false,
			Message:   "Tipo de entrada no válido",
			ErrorCode: "INVALID_ENTRY_TYPE",
		}, nil
	}
}


func (uc *ParkingUsageUsecase) processOCREntry(request *EntryRequest) (*EntryResponse, error) {
	fmt.Println("Procesando entrada OCR con patente:", request.Plate)

	if request.Plate == "" {
		return &EntryResponse{
			Success:   false,
			Message:   "Patente requerida para entrada OCR",
			ErrorCode: "PLATE_REQUIRED",
		}, nil
	}

	
	vehicle, err := uc.findVehicleByPlate(request.Plate)
	if err != nil || vehicle == nil {
		response := &EntryResponse{
			Success:   false,
			Message:   fmt.Sprintf("Vehículo con patente %s no encontrado o no autorizado", request.Plate),
			ErrorCode: "VEHICLE_NOT_FOUND",
		}
		uc.notifyEntryRejection(response, request)
		return response, nil
	}

	
	hasActiveEntry, err := uc.vehicleHasActiveEntry(vehicle.ID)
	if err != nil {
		return nil, fmt.Errorf("error al verificar ingreso activo: %w", err)
	}

	if hasActiveEntry {
		response := &EntryResponse{
			Success:   false,
			Message:   fmt.Sprintf("El vehículo %s ya tiene un ingreso activo", request.Plate),
			ErrorCode: "VEHICLE_ALREADY_ACTIVE",
		}
		uc.notifyEntryRejection(response, request)
		return response, nil
	}

	
	reservation, err := uc.findPendingReservationByVehicle(vehicle.ID)
	if err != nil {
		return nil, fmt.Errorf("error al buscar reserva: %w", err)
	}

	
	parkingUsage := &entities.ParkingUsage{
		VehicleID:   &vehicle.ID,
		OcrPlate:    request.Plate,
		QrScanned:   false,
		ManualEntry: false,
		EntryTime:   nil, 
	}

	
	if reservation != nil {
		fmt.Printf("Reserva encontrada ID: %d, usando parking ID: %d\n", reservation.ID, reservation.ParkingID)
		parkingUsage.ReservationID = &reservation.ID
		parkingUsage.ParkingID = reservation.ParkingID
	} else {
		fmt.Println("No hay reserva para OCR, se buscará parking disponible automáticamente")
		
		parkingUsage.ParkingID = 0
	}

	return uc.createParkingUsageEntry(parkingUsage, request)
}


func (uc *ParkingUsageUsecase) processQREntry(request *EntryRequest) (*EntryResponse, error) {
	if request.VehicleID == 0 {
		return &EntryResponse{
			Success:   false,
			Message:   "Vehicle ID requerido para entrada QR",
			ErrorCode: "VEHICLE_ID_REQUIRED",
		}, nil
	}

	fmt.Println("Procesando entrada QR para vehicle ID:", request.VehicleID)

	
	vehicle, err := uc.findVehicleByID(request.VehicleID)
	if err != nil || vehicle == nil {
		response := &EntryResponse{
			Success:   false,
			Message:   "Vehículo no encontrado",
			ErrorCode: "VEHICLE_NOT_FOUND",
		}
		uc.notifyEntryRejection(response, request)
		return response, nil
	}

	
	hasActiveEntry, err := uc.vehicleHasActiveEntry(request.VehicleID)
	if err != nil {
		return nil, fmt.Errorf("error al verificar ingreso activo: %w", err)
	}

	if hasActiveEntry {
		response := &EntryResponse{
			Success:   false,
			Message:   fmt.Sprintf("El vehículo %s ya tiene un ingreso activo", vehicle.Plate),
			ErrorCode: "VEHICLE_ALREADY_ACTIVE",
		}
		uc.notifyEntryRejection(response, request)
		return response, nil
	}

	
	reservation, err := uc.findPendingReservationByVehicle(request.VehicleID)
	if err != nil {
		return nil, fmt.Errorf("error al buscar reserva: %w", err)
	}

	
	parkingUsage := &entities.ParkingUsage{
		VehicleID:   &request.VehicleID,
		OcrPlate:    vehicle.Plate,
		QrScanned:   true,
		ManualEntry: false,
		EntryTime:   nil,
	}

	
	if reservation != nil {
		fmt.Printf("Reserva encontrada ID: %d, usando parking ID: %d\n", reservation.ID, reservation.ParkingID)
		parkingUsage.ReservationID = &reservation.ID
		parkingUsage.ParkingID = reservation.ParkingID
	} else {
		fmt.Println("No hay reserva, se buscará parking disponible automáticamente")
		
		parkingUsage.ParkingID = 0
	}

	return uc.createParkingUsageEntry(parkingUsage, request)
}


func (uc *ParkingUsageUsecase) processManualEntry(request *EntryRequest) (*EntryResponse, error) {
	if request.VisitorName == "" || request.Plate == "" {
		return &EntryResponse{
			Success:   false,
			Message:   "Nombre del visitante y patente son requeridos",
			ErrorCode: "VISITOR_DATA_REQUIRED",
		}, nil
	}

	
	
	vehicle, err := uc.findVehicleByPlate(request.Plate)
	if err == nil && vehicle != nil {
		
		hasActiveEntry, err := uc.vehicleHasActiveEntry(vehicle.ID)
		if err != nil {
			return nil, fmt.Errorf("error al verificar ingreso activo: %w", err)
		}

		if hasActiveEntry {
			response := &EntryResponse{
				Success:   false,
				Message:   fmt.Sprintf("El vehículo %s ya tiene un ingreso activo", request.Plate),
				ErrorCode: "VEHICLE_ALREADY_ACTIVE",
			}
			uc.notifyEntryRejection(response, request)
			return response, nil
		}
	}

	
	parkingUsage := &entities.ParkingUsage{
		ParkingID:      request.ParkingID, 
		OcrPlate:       request.Plate,
		QrScanned:      false,
		ManualEntry:    true,
		VisitorName:    request.VisitorName,
		VisitorRut:     request.VisitorRut,
		VisitorContact: request.VisitorContact,
		Zone:           request.Zone,
		RegisteredBy:   &request.RegisteredBy,
		EntryTime:      nil,
	}

	
	if vehicle != nil {
		parkingUsage.VehicleID = &vehicle.ID

		
		reservation, err := uc.findPendingReservationByVehicle(vehicle.ID)
		if err == nil && reservation != nil {
			fmt.Printf("Reserva encontrada para entrada manual ID: %d\n", reservation.ID)
			parkingUsage.ReservationID = &reservation.ID
			if parkingUsage.ParkingID == 0 {
				parkingUsage.ParkingID = reservation.ParkingID
			}
		}
	}

	return uc.createParkingUsageEntry(parkingUsage, request)
}


func (uc *ParkingUsageUsecase) createParkingUsageEntry(parkingUsage *entities.ParkingUsage, request *EntryRequest) (*EntryResponse, error) {
	
	err := uc.CreateParkingUsage(parkingUsage)
	if err != nil {
		response := &EntryResponse{
			Success:   false,
			Message:   err.Error(),
			ErrorCode: "PARKING_NOT_AVAILABLE",
		}
		uc.notifyEntryRejection(response, request)
		return response, nil
	}

	
	response := &EntryResponse{
		Success:      true,
		Message:      "Ingreso registrado exitosamente",
		ParkingUsage: parkingUsage,
	}

	
	uc.notifyEntrySuccess(response, request)

	return response, nil
}


func (uc *ParkingUsageUsecase) findVehicleByPlate(plate string) (*vehicleEntity.Vehicle, error) {
	return uc.VehicleRepository.SearchVehicleByPlate(plate)
}

func (uc *ParkingUsageUsecase) findVehicleByID(vehicleID int) (*vehicleEntity.Vehicle, error) {
	return uc.VehicleRepository.SearchVehicleByID(vehicleID)
}

func (uc *ParkingUsageUsecase) findPendingReservationByVehicle(vehicleID int) (*reservationEntity.Reservation, error) {
	reservationsEntity, err := uc.ReservationRepository.SearchPendingReservationsByVehicleID(vehicleID)
	if err != nil {
		return nil, err
	}

	if reservationsEntity == nil || len(*reservationsEntity) == 0 {
		return nil, nil
	}

	
	firstReservation := (*reservationsEntity)[0]
	return &firstReservation, nil
}


func (uc *ParkingUsageUsecase) consumeReservation(reservationID int) error {
	reservation, err := uc.ReservationRepository.SearchReservationByID(reservationID)
	if err != nil {
		return fmt.Errorf("error al buscar reserva: %w", err)
	}

	if reservation == nil {
		return fmt.Errorf("reserva no encontrada")
	}

	
	reservation.Status = "active"

	err = uc.ReservationRepository.UpdateReservationByID(reservation)
	if err != nil {
		return fmt.Errorf("error al actualizar reserva: %w", err)
	}

	return nil
}


func (uc *ParkingUsageUsecase) notifyEntryRejection(response *EntryResponse, request *EntryRequest) {
	if uc.WebSocketService != nil {
		notification := map[string]interface{}{
			"type":      "entry_rejected",
			"response":  response,
			"request":   request,
			"timestamp": time.Now(),
		}
		uc.WebSocketService.BroadcastParkingUsage(notification)
	}
}

func (uc *ParkingUsageUsecase) notifyEntrySuccess(response *EntryResponse, request *EntryRequest) {
	if uc.WebSocketService != nil {
		notification := map[string]interface{}{
			"type":      "entry_success",
			"response":  response,
			"request":   request,
			"timestamp": time.Now(),
		}
		uc.WebSocketService.BroadcastParkingUsage(notification)
	}
}


func (uc *ParkingUsageUsecase) notifyExitRejection(response *ExitResponse, request *ExitRequest) {
	if uc.WebSocketService != nil {
		notification := map[string]interface{}{
			"type":      "exit_rejected",
			"response":  response,
			"request":   request,
			"timestamp": time.Now(),
		}
		uc.WebSocketService.BroadcastParkingUsage(notification)
	}
}

func (uc *ParkingUsageUsecase) notifyExitSuccess(response *ExitResponse, request *ExitRequest) {
	if uc.WebSocketService != nil {
		notification := map[string]interface{}{
			"type":      "exit_success",
			"response":  response,
			"request":   request,
			"timestamp": time.Now(),
		}
		uc.WebSocketService.BroadcastParkingUsage(notification)
	}
}


func (uc *ParkingUsageUsecase) vehicleHasActiveEntry(vehicleID int) (bool, error) {
	activeUsages, err := uc.ParkingUsageRepository.SearchActiveParkingUsages()
	if err != nil {
		return false, err
	}

	for _, usage := range *activeUsages {
		
		if usage.VehicleID != nil && *usage.VehicleID == vehicleID && usage.ExitTime == nil {
			return true, nil
		}
	}

	return false, nil
}


func (uc *ParkingUsageUsecase) parkingHasActiveUsage(parkingID int) (bool, error) {
	activeUsages, err := uc.ParkingUsageRepository.SearchActiveParkingUsages()
	if err != nil {
		return false, err
	}

	for _, usage := range *activeUsages {
		if usage.ParkingID == parkingID && usage.ExitTime == nil {
			return true, nil
		}
	}

	return false, nil
}


func (uc *ParkingUsageUsecase) GetParkingUsagesByCustomerID(customerID string, filters map[string]interface{}) (*[]entities.ParkingUsage, error) {
	
	vehicles, err := uc.VehicleRepository.SearchVehiclesByCustomerID(customerID)
	if err != nil {
		return nil, fmt.Errorf("error al buscar vehículos del cliente: %w", err)
	}

	if len(*vehicles) == 0 {
		
		return &[]entities.ParkingUsage{}, nil
	}

	
	vehicleIDs := make([]int, len(*vehicles))
	for i, vehicle := range *vehicles {
		vehicleIDs[i] = vehicle.ID
	}

	
	result, err := uc.ParkingUsageRepository.SearchParkingUsagesByVehicleIDs(vehicleIDs, filters)
	if err != nil {
		return nil, fmt.Errorf("error al buscar registros de uso: %w", err)
	}

	return result, nil
}


func (uc *ParkingUsageUsecase) getActiveReservationForParking(parkingID int, currentTime time.Time) (*int, error) {
	reservation, err := uc.ReservationRepository.SearchActiveReservationByParkingAndTime(parkingID, currentTime)
	if err != nil {
		return nil, err
	}

	if reservation == nil {
		return nil, nil 
	}

	return &reservation.ID, nil
}

package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"

	parkingRepository "github.com/gonzalohonorato/servercorego/core/parking/domain/repositories"
	"github.com/gonzalohonorato/servercorego/core/parkingusage/application"
	"github.com/gonzalohonorato/servercorego/core/parkingusage/domain/entities"
	"github.com/gonzalohonorato/servercorego/core/parkingusage/domain/repositories"
	reservationRepository "github.com/gonzalohonorato/servercorego/core/reservation/domain/repositories"
	vehicleRepository "github.com/gonzalohonorato/servercorego/core/vehicle/domain/repositories"
	"github.com/gonzalohonorato/servercorego/core/websocket/infrastructure"
	"github.com/gorilla/mux"
)

type ParkingUsageController struct {
	ParkingUsageUsecase *application.ParkingUsageUsecase
}

func NewParkingUsageController(
	parkingUsageRepository repositories.ParkingUsageRepository,
	parkingRepository parkingRepository.ParkingRepository,
	vehicleRepository vehicleRepository.VehicleRepository,
	reservationRepository reservationRepository.ReservationRepository,
	wsService *infrastructure.WebSocketService,
) *ParkingUsageController {
	parkingUsageUseCase := application.NewParkingUsageUsecase(
		parkingUsageRepository,
		parkingRepository,
		vehicleRepository,
		reservationRepository,
		wsService,
	)

	return &ParkingUsageController{
		ParkingUsageUsecase: parkingUsageUseCase,
	}
}
func (uc *ParkingUsageController) GetParkingUsageByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	parkingUsageID := vars["id"]
	idInt, err := strconv.Atoi(parkingUsageID)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	parkingUsage, err := uc.ParkingUsageUsecase.SearchParkingUsageByID(idInt)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "ParkingUsage not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(parkingUsage)
}
func (uc *ParkingUsageController) GetParkingUsages(w http.ResponseWriter, r *http.Request) {
	parkingUsages, err := uc.ParkingUsageUsecase.SearchParkingUsages()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "ParkingUsages not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(parkingUsages)
}
func (uc *ParkingUsageController) PostRegisterExitTime(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	parkingUsageID := vars["id"]
	idInt, err := strconv.Atoi(parkingUsageID)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	exitRequest := &application.ExitRequest{
		ExitType:       "id",
		ParkingUsageID: idInt,
	}

	response, err := uc.ParkingUsageUsecase.ProcessParkingExit(exitRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error processing exit: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if !response.Success {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
func (uc *ParkingUsageController) GetActiveParkingUsages(w http.ResponseWriter, r *http.Request) {
	parkingUsages, err := uc.ParkingUsageUsecase.SearchActiveParkingUsages()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "ParkingUsages not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(parkingUsages)
}

func (uc *ParkingUsageController) GetParkingUsagesByVehicleID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vehicleID := vars["id"]
	idInt, err := strconv.Atoi(vehicleID)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	parkingUsages, err := uc.ParkingUsageUsecase.SearchParkingUsageByVehicleID(idInt)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "ParkingUsages not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(parkingUsages)
}
func (uc *ParkingUsageController) PostParkingUsage(w http.ResponseWriter, r *http.Request) {
	var newParkingUsage entities.ParkingUsage
	if err := json.NewDecoder(r.Body).Decode(&newParkingUsage); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := uc.ParkingUsageUsecase.CreateParkingUsage(&newParkingUsage); err != nil {
		http.Error(w, "Error creating parking usage", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
func (uc *ParkingUsageController) PutParkingUsage(w http.ResponseWriter, r *http.Request) {
	var newParkingUsage entities.ParkingUsage
	if err := json.NewDecoder(r.Body).Decode(&newParkingUsage); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := uc.ParkingUsageUsecase.UpdateParkingUsageById(&newParkingUsage); err != nil {
		http.Error(w, "Error update parking usage", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (uc *ParkingUsageController) DeleteParkingUsageByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	parkingUsageID := vars["id"]
	idInt, err := strconv.Atoi(parkingUsageID)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	err = uc.ParkingUsageUsecase.DeleteParkingUsageByID(idInt)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "ParkingUsage not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}

func (uc *ParkingUsageController) PostOCREntry(w http.ResponseWriter, r *http.Request) {
	var ocrRequest struct {
		Plate     string `json:"plate"`
		ParkingID int    `json:"parkingId,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&ocrRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	entryRequest := &application.EntryRequest{
		EntryType: "ocr",
		Plate:     ocrRequest.Plate,
		ParkingID: ocrRequest.ParkingID,
	}

	response, err := uc.ParkingUsageUsecase.ProcessParkingEntry(entryRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error processing entry: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if !response.Success {
		w.WriteHeader(http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(response)
}

func (uc *ParkingUsageController) PostParkingEntry(w http.ResponseWriter, r *http.Request) {
	var entryRequest application.EntryRequest
	if err := json.NewDecoder(r.Body).Decode(&entryRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	response, err := uc.ParkingUsageUsecase.ProcessParkingEntry(&entryRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error processing entry: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if !response.Success {
		w.WriteHeader(http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(response)
}

func (uc *ParkingUsageController) PostParkingExit(w http.ResponseWriter, r *http.Request) {
	var exitRequest application.ExitRequest
	if err := json.NewDecoder(r.Body).Decode(&exitRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	response, err := uc.ParkingUsageUsecase.ProcessParkingExit(&exitRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error processing exit: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if !response.Success {
		w.WriteHeader(http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(response)
}

func (uc *ParkingUsageController) GetParkingUsagesByCustomerID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID := vars["customerID"]

	filters := make(map[string]interface{})

	if parkingID := r.URL.Query().Get("parkingId"); parkingID != "" {
		if id, err := strconv.Atoi(parkingID); err == nil {
			filters["parkingId"] = id
		}
	}

	if status := r.URL.Query().Get("status"); status != "" {
		filters["status"] = status
	}

	if startDate := r.URL.Query().Get("startDate"); startDate != "" {
		filters["startDate"] = startDate
	}

	if endDate := r.URL.Query().Get("endDate"); endDate != "" {
		filters["endDate"] = endDate
	}

	result, err := uc.ParkingUsageUsecase.GetParkingUsagesByCustomerID(customerID, filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

type OCRRequest struct {
	Image     string `json:"image"`
	Timestamp string `json:"timestamp"`
}

type OCRResponse struct {
	Success      bool                   `json:"success"`
	Message      string                 `json:"message"`
	Plate        string                 `json:"plate,omitempty"`
	OcrPlate     string                 `json:"ocrPlate,omitempty"`
	Confidence   float64                `json:"confidence,omitempty"`
	ParkingUsage *entities.ParkingUsage `json:"parkingUsage,omitempty"`
	ErrorCode    string                 `json:"errorCode,omitempty"`
}

type TogetherAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (uc *ParkingUsageController) PostOCREntryWithImage(w http.ResponseWriter, r *http.Request) {
	var ocrRequest OCRRequest
	if err := json.NewDecoder(r.Body).Decode(&ocrRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if ocrRequest.Image == "" {
		http.Error(w, "Image is required", http.StatusBadRequest)
		return
	}

	plate, confidence, err := uc.processImageWithTogetherAI(ocrRequest.Image)
	if err != nil {
		response := OCRResponse{
			Success:   false,
			Message:   "Error processing image: " + err.Error(),
			ErrorCode: "OCR_PROCESSING_ERROR",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if plate == "" {
		response := OCRResponse{
			Success:   false,
			Message:   "No license plate detected in image",
			ErrorCode: "NO_PLATE_DETECTED",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	entryRequest := &application.EntryRequest{
		EntryType: "ocr",
		Plate:     plate,
		ParkingID: 0,
	}

	entryResponse, err := uc.ParkingUsageUsecase.ProcessParkingEntry(entryRequest)
	if err != nil {
		response := OCRResponse{
			Success:   false,
			Message:   "Error processing entry: " + err.Error(),
			ErrorCode: "ENTRY_PROCESSING_ERROR",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	ocrResponse := OCRResponse{
		Success:      entryResponse.Success,
		Message:      entryResponse.Message,
		Plate:        plate,
		OcrPlate:     plate,
		Confidence:   confidence,
		ParkingUsage: entryResponse.ParkingUsage,
		ErrorCode:    entryResponse.ErrorCode,
	}

	w.Header().Set("Content-Type", "application/json")
	if !ocrResponse.Success {
		w.WriteHeader(http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(ocrResponse)
}

func (uc *ParkingUsageController) PostOCRExit(w http.ResponseWriter, r *http.Request) {
	var ocrRequest OCRRequest
	if err := json.NewDecoder(r.Body).Decode(&ocrRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if ocrRequest.Image == "" {
		http.Error(w, "Image is required", http.StatusBadRequest)
		return
	}

	plate, confidence, err := uc.processImageWithTogetherAI(ocrRequest.Image)
	if err != nil {
		response := OCRResponse{
			Success:   false,
			Message:   "Error processing image: " + err.Error(),
			ErrorCode: "OCR_PROCESSING_ERROR",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if plate == "" {
		response := OCRResponse{
			Success:   false,
			Message:   "No license plate detected in image",
			ErrorCode: "NO_PLATE_DETECTED",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	exitRequest := &application.ExitRequest{
		ExitType: "plate",
		Plate:    plate,
	}

	exitResponse, err := uc.ParkingUsageUsecase.ProcessParkingExit(exitRequest)
	if err != nil {
		response := OCRResponse{
			Success:   false,
			Message:   "Error processing exit: " + err.Error(),
			ErrorCode: "EXIT_PROCESSING_ERROR",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	ocrResponse := OCRResponse{
		Success:      exitResponse.Success,
		Message:      exitResponse.Message,
		Plate:        plate,
		OcrPlate:     plate,
		Confidence:   confidence,
		ParkingUsage: exitResponse.ParkingUsage,
		ErrorCode:    exitResponse.ErrorCode,
	}

	w.Header().Set("Content-Type", "application/json")
	if !ocrResponse.Success {
		w.WriteHeader(http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(ocrResponse)
}

func (uc *ParkingUsageController) processImageWithTogetherAI(base64Image string) (string, float64, error) {

	imageDataURL := "data:image/jpeg;base64," + base64Image

	payload := map[string]interface{}{
		"model": "meta-llama/Llama-Vision-Free",
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": `Analiza esta imagen y extrae ÚNICAMENTE la patente del vehículo. 
						Responde con el siguiente formato JSON exacto:
						{"plate": "ABC123", "confidence": 0.95}
						
						Si no puedes detectar una patente claramente, responde:
						{"plate": "", "confidence": 0.0}
						
						La patente debe estar en formato chileno (6-7 caracteres alfanuméricos).
						NO agregues explicaciones adicionales, solo el JSON.`,
					},
					{
						"type": "image_url",
						"image_url": map[string]string{
							"url": imageDataURL,
						},
					},
				},
			},
		},
		"max_tokens":  100,
		"temperature": 0.1,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", 0, fmt.Errorf("error marshaling payload: %w", err)
	}

	req, err := http.NewRequest("POST", "https:")
	if err != nil {
		return "", 0, fmt.Errorf("error creating request: %w", err)
	}

	apiKey := os.Getenv("TOGETHER_API_KEY")
	if apiKey == "" {
		return "", 0, fmt.Errorf("TOGETHER_API_KEY environment variable not set")
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("error sending request to Together.AI: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("Together.AI API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var togetherResp TogetherAIResponse
	if err := json.Unmarshal(respBody, &togetherResp); err != nil {
		return "", 0, fmt.Errorf("error parsing Together.AI response: %w", err)
	}

	if len(togetherResp.Choices) == 0 {
		return "", 0, fmt.Errorf("no choices in Together.AI response")
	}

	content := togetherResp.Choices[0].Message.Content
	fmt.Println("Contenido de la respuesta:", content)

	var plateResult struct {
		Plate      string  `json:"plate"`
		Confidence float64 `json:"confidence"`
	}

	if err := json.Unmarshal([]byte(content), &plateResult); err != nil {

		plate := uc.extractPlateFromText(content)
		if plate != "" {
			return plate, 0.8, nil
		}
		return "", 0, fmt.Errorf("could not parse plate from response: %s", content)
	}

	if plateResult.Plate != "" && uc.isValidChileanPlate(plateResult.Plate) {
		return plateResult.Plate, plateResult.Confidence, nil
	}

	return "", 0, nil
}

func (uc *ParkingUsageController) extractPlateFromText(text string) string {

	patterns := []string{
		`[A-Z]{4}[0-9]{2}`,
		`[A-Z]{2}[0-9]{4}`,
		`[A-Z]{3}[0-9]{3}`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if match := re.FindString(text); match != "" {
			return match
		}
	}

	return ""
}

func (uc *ParkingUsageController) isValidChileanPlate(plate string) bool {

	if len(plate) < 6 || len(plate) > 7 {
		return false
	}

	patterns := []string{
		`^[A-Z]{4}[0-9]{2}$`,
		`^[A-Z]{2}[0-9]{4}$`,
		`^[A-Z]{3}[0-9]{3}$`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if re.MatchString(plate) {
			return true
		}
	}

	return false
}

package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"firebase.google.com/go/v4/auth" 
	"github.com/gonzalohonorato/servercorego/core/user/application"
	"github.com/gonzalohonorato/servercorego/core/user/domain/entities"
	"github.com/gonzalohonorato/servercorego/core/user/domain/repositories"
	"github.com/gorilla/mux"
)

type UserController struct {
	UserUsecase *application.UserUsecase
}

func NewUserController(userRepository repositories.UserRepository) *UserController {
	userUseCase := application.NewUserUsecase(userRepository)
	return &UserController{
		UserUsecase: userUseCase,
	}
}


func (uc *UserController) SetFirebaseAuth(firebaseAuth *auth.Client) {
	uc.UserUsecase.SetFirebaseAuth(firebaseAuth)
}


type CreateCustomerResponse struct {
	Success   bool           `json:"success"`
	Message   string         `json:"message"`
	User      *entities.User `json:"user,omitempty"`
	Vehicle   interface{}    `json:"vehicle,omitempty"`
	ErrorCode string         `json:"errorCode,omitempty"`
}

func (uc *UserController) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var req application.CreateCustomerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CreateCustomerResponse{
			Success:   false,
			Message:   "Datos inválidos: " + err.Error(),
			ErrorCode: "INVALID_DATA",
		})
		return
	}

	user, vehicle, err := uc.UserUsecase.CreateCustomerWithFirebase(req)

	if err != nil {
		log.Printf("Error creating customer: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(CreateCustomerResponse{
			Success:   false,
			Message:   "Error al crear cliente: " + err.Error(),
			ErrorCode: "CREATION_ERROR",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateCustomerResponse{
		Success: true,
		Message: "Cliente creado exitosamente",
		User:    user,
		Vehicle: vehicle,
	})
}


type CreateEmployeeResponse struct {
	Success   bool           `json:"success"`
	Message   string         `json:"message"`
	User      *entities.User `json:"user,omitempty"`
	ErrorCode string         `json:"errorCode,omitempty"`
}

func (uc *UserController) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var req application.CreateEmployeeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CreateEmployeeResponse{
			Success:   false,
			Message:   "Datos inválidos: " + err.Error(),
			ErrorCode: "INVALID_DATA",
		})
		return
	}

	user, err := uc.UserUsecase.CreateEmployeeWithFirebase(req)

	if err != nil {
		log.Printf("Error creating employee: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(CreateEmployeeResponse{
			Success:   false,
			Message:   "Error al crear empleado: " + err.Error(),
			ErrorCode: "CREATION_ERROR",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateEmployeeResponse{
		Success: true,
		Message: "Empleado creado exitosamente",
		User:    user,
	})
}


func (uc *UserController) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	user, err := uc.UserUsecase.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser entities.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := uc.UserUsecase.CreateUser(&newUser); err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user entities.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := uc.UserUsecase.UpdateUser(&user); err != nil {
		http.Error(w, "Error update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (uc *UserController) SearchUsers(w http.ResponseWriter, r *http.Request) {
	users, err := uc.UserUsecase.SearchUsers()
	if err != nil {
		http.Error(w, "Error searching users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (uc *UserController) SearchUsersByType(w http.ResponseWriter, r *http.Request) {
	userType := r.URL.Query().Get("type")
	if userType == "" {
		userType = "customer"
	}

	users, err := uc.UserUsecase.SearchUsersByType(userType)
	if err != nil {
		http.Error(w, "Error searching users by type", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

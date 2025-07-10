package application

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/gonzalohonorato/servercorego/core/user/domain/entities"
	"github.com/gonzalohonorato/servercorego/core/user/domain/repositories"
)

type UserUsecase struct {
	UserRepository repositories.UserRepository
	FirebaseAuth   *auth.Client
}

func NewUserUsecase(userRepo repositories.UserRepository) *UserUsecase {
	return &UserUsecase{
		UserRepository: userRepo,
	}
}

func (uc *UserUsecase) SetFirebaseAuth(firebaseAuth *auth.Client) {
	uc.FirebaseAuth = firebaseAuth
}


type CreateCustomerRequest struct {
	Name                   string       `json:"name"`
	Email                  string       `json:"email"`
	Rut                    string       `json:"rut"`
	CustomerType           *string      `json:"customerType"`
	Password               string       `json:"password"`
	GenerateRandomPassword bool         `json:"generateRandomPassword"`
	SendEmailPassword      bool         `json:"sendEmailPassword"`
	Vehicle                *VehicleData `json:"vehicle,omitempty"`
}

type VehicleData struct {
	Plate string `json:"plate"`
	Brand string `json:"brand"`
	Model string `json:"model"`
}

func (uc *UserUsecase) CreateCustomerWithFirebase(req CreateCustomerRequest) (*entities.User, interface{}, error) {
	ctx := context.Background()

	
	finalPassword := req.Password
	if req.GenerateRandomPassword {
		generatedPassword, err := generateRandomPassword()
		if err != nil {
			return nil, nil, fmt.Errorf("error generating password: %w", err)
		}
		finalPassword = generatedPassword
	}

	
	if uc.FirebaseAuth == nil {
		return nil, nil, fmt.Errorf("Firebase Auth not initialized")
	}

	firebaseUser := &auth.UserToCreate{}
	firebaseUser.Email(req.Email)
	firebaseUser.Password(finalPassword)
	firebaseUser.DisplayName(req.Name)

	createdFirebaseUser, err := uc.FirebaseAuth.CreateUser(ctx, firebaseUser)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating Firebase user: %w", err)
	}

	log.Printf("Firebase user created with UID: %s", createdFirebaseUser.UID)

	
	user := &entities.User{
		ID:           createdFirebaseUser.UID,
		Name:         req.Name,
		Email:        req.Email,
		Rut:          req.Rut,
		Uid:          createdFirebaseUser.UID,
		Type:         "customer",
		CustomerType: req.CustomerType,
		CreatedAt:    time.Now(),
	}

	
	err = uc.UserRepository.CreateUser(user)
	if err != nil {
		
		uc.FirebaseAuth.DeleteUser(ctx, createdFirebaseUser.UID)
		return nil, nil, fmt.Errorf("error creating user in database: %w", err)
	}

	
	var createdVehicle interface{}
	if req.Vehicle != nil && req.Vehicle.Plate != "" {
		log.Printf("Vehicle creation would be implemented here for user %s", user.ID)

		createdVehicle = map[string]interface{}{
			"plate":      req.Vehicle.Plate,
			"brand":      req.Vehicle.Brand,
			"model":      req.Vehicle.Model,
			"customerId": user.ID,
		}
	}

	
	if req.SendEmailPassword {
		err = uc.sendPasswordEmail(req.Email, req.Name, finalPassword)
		if err != nil {
			log.Printf("Warning: User created but email sending failed: %v", err)
		}
	}

	return user, createdVehicle, nil
}


type CreateEmployeeRequest struct {
	Name                   string  `json:"name"`
	Email                  string  `json:"email"`
	Rut                    string  `json:"rut"`
	EmployeeRole           *string `json:"employeeRole"`
	Password               string  `json:"password"`
	GenerateRandomPassword bool    `json:"generateRandomPassword"`
	SendEmailPassword      bool    `json:"sendEmailPassword"`
}

func (uc *UserUsecase) CreateEmployeeWithFirebase(req CreateEmployeeRequest) (*entities.User, error) {
	ctx := context.Background()

	
	finalPassword := req.Password
	if req.GenerateRandomPassword {
		generatedPassword, err := generateRandomPassword()
		if err != nil {
			return nil, fmt.Errorf("error generating password: %w", err)
		}
		finalPassword = generatedPassword
	}

	
	if uc.FirebaseAuth == nil {
		return nil, fmt.Errorf("Firebase Auth not initialized")
	}

	firebaseUser := &auth.UserToCreate{}
	firebaseUser.Email(req.Email)
	firebaseUser.Password(finalPassword)
	firebaseUser.DisplayName(req.Name)

	createdFirebaseUser, err := uc.FirebaseAuth.CreateUser(ctx, firebaseUser)
	if err != nil {
		return nil, fmt.Errorf("error creating Firebase user: %w", err)
	}

	log.Printf("Firebase user created with UID: %s", createdFirebaseUser.UID)

	
	user := &entities.User{
		ID:           createdFirebaseUser.UID, 
		Uid:          createdFirebaseUser.UID, 
		Name:         req.Name,
		Email:        req.Email,
		Rut:          req.Rut,
		Type:         "employee",
		EmployeeRole: req.EmployeeRole,
		CreatedAt:    time.Now(),
		
	}

	
	err = uc.UserRepository.CreateUser(user)
	if err != nil {
		
		uc.FirebaseAuth.DeleteUser(ctx, createdFirebaseUser.UID)
		return nil, fmt.Errorf("error creating user in database: %w", err)
	}

	log.Printf("Employee created in PostgreSQL with ID: %s (Firebase UID)", user.ID)

	
	if req.SendEmailPassword {
		err = uc.sendPasswordEmail(req.Email, req.Name, finalPassword)
		if err != nil {
			log.Printf("Warning: Employee created but email sending failed: %v", err)
		}
	}

	return user, nil
}


func (uc *UserUsecase) GetUserByID(id string) (*entities.User, error) {
	return uc.UserRepository.SearchUserByID(id)
}

func (uc *UserUsecase) CreateUser(user *entities.User) error {
	return uc.UserRepository.CreateUser(user)
}

func (uc *UserUsecase) UpdateUser(user *entities.User) error {
	return uc.UserRepository.UpdateUser(user)
}

func (uc *UserUsecase) SearchUsers() (*entities.Users, error) {
	return uc.UserRepository.SearchUsers()
}

func (uc *UserUsecase) SearchUsersByType(userType string) (*entities.Users, error) {
	return uc.UserRepository.SearchUsersByType(userType)
}

func generateRandomPassword() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	const length = 12

	password := make([]byte, length)
	for i := range password {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[randomIndex.Int64()]
	}

	return string(password), nil
}

func (uc *UserUsecase) sendPasswordEmail(email, name, password string) error {
	log.Printf("Would send password email to %s for user %s with password %s", email, name, password)
	return nil
}

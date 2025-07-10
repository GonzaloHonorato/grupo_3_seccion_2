package entities

import (
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Rut       string    `json:"rut"`
	Uid       string    `json:"uid"`
	Type      string    `json:"type"` 
	CreatedAt time.Time `json:"createdAt"`

	
	CustomerType *string `json:"customerType,omitempty"`
	EmployeeRole *string `json:"employeeRole,omitempty"`
}

type Users []User

package routes

import (
	"log"

	"github.com/gonzalohonorato/servercorego/config/injector"
	"github.com/gonzalohonorato/servercorego/core/user/infrastructure/rest/controllers"
	"github.com/gorilla/mux"
)

func UserRoutes(router *mux.Router, container *injector.Container) {
	
	controller := controllers.NewUserController(container.ProvideUserRepository())

	
	firebaseAuth, err := container.GetFirebaseAuth()
	if err != nil {
		log.Printf("Warning: Firebase Auth not initialized: %v", err)
	} else {
		controller.SetFirebaseAuth(firebaseAuth)
	}

	
	router.HandleFunc("/users", controller.SearchUsersByType).Methods("GET")
	router.HandleFunc("/users/{id}", controller.GetUserByID).Methods("GET")
	router.HandleFunc("/users", controller.CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", controller.UpdateUser).Methods("PUT")

	
	router.HandleFunc("/customers", controller.CreateCustomer).Methods("POST")
	router.HandleFunc("/employees", controller.CreateEmployee).Methods("POST") 
}

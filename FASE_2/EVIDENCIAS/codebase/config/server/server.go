package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gonzalohonorato/servercorego/config/injector"
	feedbackRoutes "github.com/gonzalohonorato/servercorego/core/feedback/infrastructure/rest/routes"
	notificationtemplateRoutes "github.com/gonzalohonorato/servercorego/core/notificationtemplate/infrastructure/rest/routes"
	parkingRoutes "github.com/gonzalohonorato/servercorego/core/parking/infrastructure/rest/routes"
	parkingUsageRoutes "github.com/gonzalohonorato/servercorego/core/parkingusage/infrastructure/rest/routes"
	reservationRoutes "github.com/gonzalohonorato/servercorego/core/reservation/infrastructure/rest/routes"
	userRoutes "github.com/gonzalohonorato/servercorego/core/user/infrastructure/rest/routes"
	usernotificationRoutes "github.com/gonzalohonorato/servercorego/core/usernotification/infrastructure/rest/routes"
	vehicleRoutes "github.com/gonzalohonorato/servercorego/core/vehicle/infrastructure/rest/routes"
	websocketRoutes "github.com/gonzalohonorato/servercorego/core/websocket/infrastructure/rest/routes"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func RunServer() {
	log.Println("Iniciando servidor...")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("Variable de ambiente PORT no seteada")
	}

	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	
	container := injector.NewContainer(ctx)

	
	defer container.CloseTimescaleDB()

	
	router := mux.NewRouter().StrictSlash(true)

	
	

	
	

	
	
	userRoutes.UserRoutes(router, container)
	parkingRoutes.ParkingRoutes(router, container)
	vehicleRoutes.VehicleRoutes(router, container)
	parkingUsageRoutes.ParkingUsageRoutes(router, container)
	reservationRoutes.ReservationRoutes(router, container)
	feedbackRoutes.FeedbackRoutes(router, container)
	usernotificationRoutes.UserNotificationRoutes(router, container)
	notificationtemplateRoutes.NotificationTemplateRoutes(router, container)
	wsService := container.ProvideWebSocketService()

	
	websocketRoutes.WebSocketRoutes(router, wsService)

	
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	
	handler := corsHandler.Handler(router)

	log.Printf("Servidor escuchando en puerto %s", port)

	
	reservationScheduler := container.ProvideReservationScheduler()
	reservationScheduler.Start()

	
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Deteniendo servicios...")
		reservationScheduler.Stop()
		container.CloseTimescaleDB()
		os.Exit(0)
	}()

	panic(http.ListenAndServe(":"+port, handler))
}

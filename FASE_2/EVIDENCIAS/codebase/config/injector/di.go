package injector

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/v4/auth"
	"github.com/gonzalohonorato/servercorego/config"
	feedbackPersistence "github.com/gonzalohonorato/servercorego/core/feedback/infrastructure/persistence"
	notificationtemplatePersistence "github.com/gonzalohonorato/servercorego/core/notificationtemplate/infrastructure/persistence"
	parkingPersistence "github.com/gonzalohonorato/servercorego/core/parking/infrastructure/persistence"
	parkingUsagePersistence "github.com/gonzalohonorato/servercorego/core/parkingusage/infrastructure/persistence"
	"github.com/gonzalohonorato/servercorego/core/reservation/application"
	reservation "github.com/gonzalohonorato/servercorego/core/reservation/application"
	reservationPersistence "github.com/gonzalohonorato/servercorego/core/reservation/infrastructure/persistence"
	userPersistence "github.com/gonzalohonorato/servercorego/core/user/infrastructure/persistence"
	usernotification "github.com/gonzalohonorato/servercorego/core/usernotification/infrastructure/persistence"
	vehiclePersistence "github.com/gonzalohonorato/servercorego/core/vehicle/infrastructure/persistence"
	"github.com/gonzalohonorato/servercorego/core/websocket/infrastructure"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	Ctx context.Context

	dbPool *pgxpool.Pool
	dbOnce sync.Once

	
	FirestoreClient *firestore.Client
	FirebaseAuth    *auth.Client
	FirebaseOnce    sync.Once

	wsService *infrastructure.WebSocketService
	wsOnce    sync.Once

	reservationScheduler *application.ReservationScheduler
	schedulerOnce        sync.Once
}

func NewContainer(ctx context.Context) *Container {
	return &Container{Ctx: ctx}
}

func (c *Container) InitTimescaleDB() (*pgxpool.Pool, error) {
	var err error

	c.dbOnce.Do(func() {
		connStr := os.Getenv("DATABASE_URL")
		if connStr == "" {
			err = fmt.Errorf("Variable de ambiente DATABASE_URL no seteada")
			return
		}

		pool, e := pgxpool.New(c.Ctx, connStr)
		if e != nil {
			err = fmt.Errorf("Error al conectar a TimescaleDB: %v", e)
			return
		}

		c.dbPool = pool
	})

	return c.dbPool, err
}

func (c *Container) CloseTimescaleDB() {
	if c.dbPool != nil {
		c.dbPool.Close()
	}
}


func (c *Container) InitFirebase() (*config.FirebaseClients, error) {
	var err error
	var clients *config.FirebaseClients

	c.FirebaseOnce.Do(func() {
		firebaseClients, firebaseErr := config.InitFirebase()
		if firebaseErr != nil {
			err = firebaseErr
			return
		}

		c.FirestoreClient = firebaseClients.Firestore
		c.FirebaseAuth = firebaseClients.Auth
		clients = firebaseClients
	})

	if clients == nil && err == nil {
		
		clients = &config.FirebaseClients{
			Firestore: c.FirestoreClient,
			Auth:      c.FirebaseAuth,
		}
	}

	return clients, err
}


func (c *Container) InitFirestore() (*firestore.Client, error) {
	clients, err := c.InitFirebase()
	if err != nil {
		return nil, err
	}
	return clients.Firestore, nil
}


func (c *Container) GetFirebaseAuth() (*auth.Client, error) {
	clients, err := c.InitFirebase()
	if err != nil {
		return nil, err
	}
	return clients.Auth, nil
}

func (c *Container) CloseFirestore() {
	if c.FirestoreClient != nil {
		_ = c.FirestoreClient.Close()
	}
}


func (c *Container) ProvideParkingRepository() *parkingPersistence.TimescaleParkingRepository {
	pool, err := c.InitTimescaleDB()
	if err != nil {
		log.Fatalf("Error al inicializar TimescaleDB: %v", err)
	}
	return parkingPersistence.NewTimescaleDBRepository(pool)
}

func (c *Container) ProvideVehicleRepository() *vehiclePersistence.TimescaleVehicleRepository {
	pool, err := c.InitTimescaleDB()
	if err != nil {
		log.Fatalf("Error al inicializar TimescaleDB: %v", err)
	}
	return vehiclePersistence.NewTimescaleVehicleRepository(pool)
}

func (c *Container) ProvideReservationRepository() *reservationPersistence.TimescaleReservationRepository {
	pool, err := c.InitTimescaleDB()
	if err != nil {
		log.Fatalf("Error al inicializar TimescaleDB: %v", err)
	}
	return reservationPersistence.NewTimescaleDBRepository(pool)
}

func (c *Container) ProvideParkingUsageRepository() *parkingUsagePersistence.TimescaleParkingUsageRepository {
	pool, err := c.InitTimescaleDB()
	if err != nil {
		log.Fatalf("Error al inicializar TimescaleDB: %v", err)
	}
	return parkingUsagePersistence.NewTimescaleDBRepository(pool)
}

func (c *Container) ProvideUserRepository() *userPersistence.TimescaleUserRepository {
	pool, err := c.InitTimescaleDB()
	if err != nil {
		log.Fatalf("Error al inicializar TimescaleDB: %v", err)
	}
	return userPersistence.NewTimescaleDBRepository(pool)
}

func (c *Container) ProvideFeedbackRepository() *feedbackPersistence.TimescaleFeedbackRepository {
	pool, err := c.InitTimescaleDB()
	if err != nil {
		log.Fatalf("Error al inicializar TimescaleDB: %v", err)
	}
	return feedbackPersistence.NewTimescaleDBRepository(pool)
}

func (c *Container) ProvideUserNotificationRepository() *usernotification.TimescaleUserNotificationRepository {
	pool, err := c.InitTimescaleDB()
	if err != nil {
		log.Fatalf("Error al inicializar TimescaleDB: %v", err)
	}
	return usernotification.NewTimescaleUserNotificationRepository(pool)
}

func (c *Container) ProvideNotificationTemplateRepository() *notificationtemplatePersistence.TimescaleNotificationTemplateRepository {
	pool, err := c.InitTimescaleDB()
	if err != nil {
		log.Fatalf("Error al inicializar TimescaleDB: %v", err)
	}
	return notificationtemplatePersistence.NewTimescaleNotificationTemplateRepository(pool)
}

func (c *Container) ProvideWebSocketService() *infrastructure.WebSocketService {
	c.wsOnce.Do(func() {
		c.wsService = infrastructure.NewWebSocketService()

		go c.wsService.Start()
	})

	return c.wsService
}




func (c *Container) ProvideReservationScheduler() *application.ReservationScheduler {
	c.schedulerOnce.Do(func() {
		reservationRepo := c.ProvideReservationRepository()
		parkingRepo := c.ProvideParkingRepository()

		
		reservationUsecase := reservation.NewReservationUsecase(reservationRepo, parkingRepo)

		
		c.reservationScheduler = reservation.NewReservationScheduler(reservationUsecase, parkingRepo)
	})

	return c.reservationScheduler
}

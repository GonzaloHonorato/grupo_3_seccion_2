package config

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)


type FirebaseClients struct {
	Firestore *firestore.Client
	Auth      *auth.Client
}

func InitFirebase() (*FirebaseClients, error) {
	ctx := context.Background()

	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	if projectID == "" {
		log.Fatalf("FIREBASE_PROJECT_ID environment variable not set")
	}

	
	credentialsMap := map[string]string{
		"type":                        os.Getenv("FIREBASE_TYPE"),
		"project_id":                  os.Getenv("FIREBASE_PROJECT_ID"),
		"private_key_id":              os.Getenv("FIREBASE_PRIVATE_KEY_ID"),
		"private_key":                 os.Getenv("FIREBASE_PRIVATE_KEY"),
		"client_email":                os.Getenv("FIREBASE_CLIENT_EMAIL"),
		"client_id":                   os.Getenv("FIREBASE_CLIENT_ID"),
		"auth_uri":                    os.Getenv("FIREBASE_AUTH_URI"),
		"token_uri":                   os.Getenv("FIREBASE_TOKEN_URI"),
		"auth_provider_x509_cert_url": os.Getenv("FIREBASE_AUTH_PROVIDER_X509_CERT_URL"),
		"client_x509_cert_url":        os.Getenv("FIREBASE_CLIENT_X509_CERT_URL"),
		"universe_domain":             os.Getenv("FIREBASE_UNIVERSE_DOMAIN"),
	}

	credentialsJSON, err := json.Marshal(credentialsMap)
	if err != nil {
		log.Fatalf("Error marshaling Firebase credentials: %v", err)
	}

	firebaseConfig := &firebase.Config{
		ProjectID: projectID,
	}

	opt := option.WithCredentialsJSON(credentialsJSON)

	firebaseApp, err := firebase.NewApp(ctx, firebaseConfig, opt)
	if err != nil {
		log.Fatalf("Error initializing app: %v", err)
	}

	
	firestoreClient, err := firebaseApp.Firestore(ctx)
	if err != nil {
		log.Fatalf("Error getting Firestore client: %v", err)
	}

	
	authClient, err := firebaseApp.Auth(ctx)
	if err != nil {
		log.Fatalf("Error getting Auth client: %v", err)
	}

	log.Println("Firebase connection successful (Firestore + Auth)")

	return &FirebaseClients{
		Firestore: firestoreClient,
		Auth:      authClient,
	}, nil
}


func InitFirestore() (*firestore.Client, error) {
	clients, err := InitFirebase()
	if err != nil {
		return nil, err
	}
	return clients.Firestore, nil
}

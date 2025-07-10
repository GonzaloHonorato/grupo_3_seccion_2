package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

type TimescaleDB struct {
	Ctx  context.Context
	Conn *pgx.Conn
}


func InitTimescaleDB() (*TimescaleDB, error) {
	connStr := os.Getenv("DATABASE_URL")
	fmt.Println("Connecting to TimescaleDB with connection", connStr)
	if connStr == "" {
		log.Fatalf("DATABASE_URL environment variable not set")
	}
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatalf("Unable to connect to TimescaleDB: %v", err)
	}
	return &TimescaleDB{
		Ctx:  ctx,
		Conn: conn,
	}, nil
}

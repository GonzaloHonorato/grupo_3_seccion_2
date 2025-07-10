package main

import (
	"log"

	"github.com/gonzalohonorato/servercorego/config/server"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No se pudo cargar el archivo .env, usando variables de entorno del sistema")
	}
	server.RunServer()
}

package main

import (
	"auth-service/internal/app"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("can not read .env file: %v", err)
	}
}

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatalf("error creating an application instance: %s\n", err.Error())
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("application startup error: %s\n", err.Error())
	}
}

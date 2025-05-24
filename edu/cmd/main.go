package main

import (
	"course2/internal/app"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Предупреждение: Не удалось прочитать файл .env: %v\n", err)
	}
}

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatalf("Ошибка создания экземпляра приложения: %s\n", err.Error())
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("Ошибка запуска приложения: %s\n", err.Error())
	}
}

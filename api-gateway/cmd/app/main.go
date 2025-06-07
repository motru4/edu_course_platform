package main

import (
	"api-gateway/internal/app"
	"fmt"
	"os"
)

func main() {
	// Создаем приложение
	application, err := app.New()
	if err != nil {
		fmt.Printf("Ошибка при создании приложения: %v\n", err)
		os.Exit(1)
	}

	// Запускаем сервер
	fmt.Println("Запуск API Gateway...")
	if err := application.Run(); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %v\n", err)
		os.Exit(1)
	}
}

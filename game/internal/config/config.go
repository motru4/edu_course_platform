package config

import (
	"os"
	"strconv"
)

// Config представляет конфигурацию приложения
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
}

// ServerConfig представляет конфигурацию HTTP сервера
type ServerConfig struct {
	Port string
}

// DatabaseConfig представляет конфигурацию базы данных
type DatabaseConfig struct {
	URL string
}

// AuthConfig представляет конфигурацию сервиса авторизации
type AuthConfig struct {
	GRPCAddress string
}

// New создает новый экземпляр Config с настройками из переменных окружения
func New() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8083"),
		},
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/eduplatform?sslmode=disable"),
		},
		Auth: AuthConfig{
			GRPCAddress: getEnv("AUTH_GRPC_ADDRESS", "localhost:9090"),
		},
	}, nil
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvBool получает булево значение переменной окружения или возвращает значение по умолчанию
func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolValue
}

// getEnvInt получает целочисленное значение переменной окружения или возвращает значение по умолчанию
func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

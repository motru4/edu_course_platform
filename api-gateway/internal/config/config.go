package config

import (
	"os"
	"strings"
)

// Config содержит настройки API-шлюза
type Config struct {
	Server   ServerConfig
	Services ServicesConfig
}

// ServerConfig содержит настройки HTTP-сервера
type ServerConfig struct {
	Port string
}

// ServicesConfig содержит настройки сервисов
type ServicesConfig struct {
	AuthService     ServiceConfig
	AuthGRPCService ServiceConfig
	EduService      ServiceConfig
	GameService     ServiceConfig
}

// ServiceConfig содержит настройки для сервиса
type ServiceConfig struct {
	URL string
}

// New создает новую конфигурацию с значениями по умолчанию или из переменных окружения
func New() *Config {
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("API_GATEWAY_PORT", "8090"),
		},
		Services: ServicesConfig{
			AuthService: ServiceConfig{
				URL: getEnv("AUTH_SERVICE_URL", "http://auth-service:8080"),
			},
			AuthGRPCService: ServiceConfig{
				URL: getEnv("AUTH_GRPC_SERVICE_URL", "auth-service:9090"),
			},
			EduService: ServiceConfig{
				URL: getEnv("EDU_SERVICE_URL", "http://edu-service:8081"),
			},
			GameService: ServiceConfig{
				URL: getEnv("GAME_SERVICE_URL", "http://game-service:8083"),
			},
		},
	}

	// Проверяем и нормализуем URL для HTTP сервисов
	cfg.Services.AuthService.URL = normalizeURL(cfg.Services.AuthService.URL)
	cfg.Services.EduService.URL = normalizeURL(cfg.Services.EduService.URL)
	cfg.Services.GameService.URL = normalizeURL(cfg.Services.GameService.URL)

	return cfg
}

// normalizeURL убеждается, что URL содержит схему
func normalizeURL(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "http://" + url
	}
	return url
}

// getEnv получает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

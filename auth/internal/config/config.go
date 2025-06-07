package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Token     TokenConfig
	SMTP      SMTPConfig
	OAuth     OAuthConfig
	RateLimit RateLimitConfig
	Security  SecurityConfig
}

type ServerConfig struct {
	Port     string
	GRPCPort string
}

type DatabaseConfig struct {
	URL string
}

type TokenConfig struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
	PepperStr  string
}

type SMTPConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string
	FromEmail string
}

type OAuthConfig struct {
	GoogleClientID     string
	GoogleClientSecret string
}

type RateLimitConfig struct {
	Requests int
	Period   time.Duration
}

type SecurityConfig struct {
	PasswordPepper string `env:"PASSWORD_PEPPER" envDefault:"your-default-pepper-key-replace-in-production"`
}

func New() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Port:     getEnvOrDefault("SERVER_PORT", "8080"),
			GRPCPort: getEnvOrDefault("GRPC_PORT", "9090"),
		},
		Database: DatabaseConfig{
			URL: os.Getenv("DB_URL"),
		},
		Token: TokenConfig{
			Secret:     os.Getenv("JWT_SECRET"),
			AccessTTL:  time.Hour,
			RefreshTTL: time.Hour * 168, // 7 days
			PepperStr:  os.Getenv("PEPPER_STRING"),
		},
		SMTP: SMTPConfig{
			Host: getEnvOrDefault("SMTP_HOST", "smtp.gmail.com"),
			Port: func(port string) int {
				smtpPort, _ := strconv.Atoi(getEnvOrDefault("SMTP_PORT", "587"))
				return smtpPort
			}(getEnvOrDefault("SMTP_PORT", "587")),
			Username:  os.Getenv("SMTP_USERNAME"),
			Password:  os.Getenv("SMTP_PASSWORD"),
			FromEmail: os.Getenv("SMTP_FROM_EMAIL"),
		},
		OAuth: OAuthConfig{
			GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			GoogleClientSecret: os.Getenv("GOOGLE_SECRET"),
		},
		RateLimit: RateLimitConfig{
			Requests: 50,
			Period:   time.Minute,
		},
		Security: SecurityConfig{
			PasswordPepper: getEnvOrDefault("PASSWORD_PEPPER", "your-default-pepper-key-replace-in-production"),
		},
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

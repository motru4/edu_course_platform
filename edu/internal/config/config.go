package config

import (
	"os"
)

type Config struct {
	Database struct {
		URL string
	}
	Server struct {
		Port string
	}
	Auth struct {
		GRPCAddress string
	}
}

func New() (*Config, error) {
	cfg := &Config{}

	// Database configuration
	cfg.Database.URL = os.Getenv("DB_URL")

	// Server configuration
	cfg.Server.Port = os.Getenv("SERVER_PORT")
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8081" // default port
	}

	// Auth configuration
	cfg.Auth.GRPCAddress = os.Getenv("AUTH_GRPC_ADDRESS")
	if cfg.Auth.GRPCAddress == "" {
		cfg.Auth.GRPCAddress = "localhost:50051" // default address
	}

	return cfg, nil
}

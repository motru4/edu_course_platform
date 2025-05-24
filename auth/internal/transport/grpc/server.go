package grpc

import (
	"auth-service/internal/services"

	"google.golang.org/grpc"
)

func NewServer(authService *services.AuthService) *grpc.Server {
	grpcServer := grpc.NewServer()
	RegisterAuthServiceServer(grpcServer, NewAuthServer(authService))
	return grpcServer
}

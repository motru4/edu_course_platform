package grpc

import (
	"auth-service/internal/services"
	"context"
)

type AuthServer struct {
	UnimplementedAuthServiceServer
	authService *services.AuthService
}

func NewAuthServer(authService *services.AuthService) *AuthServer {
	return &AuthServer{
		authService: authService,
	}
}

func (s *AuthServer) CheckAccess(ctx context.Context, req *CheckAccessRequest) (*CheckAccessResponse, error) {
	claims, err := s.authService.ValidateToken(req.Token)
	if err != nil {
		return &CheckAccessResponse{
			Allowed: false,
			UserId:  "",
			Error:   err.Error(),
		}, nil
	}

	// Check if user has any of the required roles
	if len(req.RequiredRoles) > 0 {
		hasRequiredRole := false
		userRole := string(claims.Role)

		for _, requiredRole := range req.RequiredRoles {
			if userRole == requiredRole {
				hasRequiredRole = true
				break
			}
		}

		if !hasRequiredRole {
			return &CheckAccessResponse{
				Allowed: false,
				UserId:  claims.UserID.String(),
				Error:   "insufficient permissions",
			}, nil
		}
	}

	return &CheckAccessResponse{
		Allowed: true,
		UserId:  claims.UserID.String(),
		Error:   "",
	}, nil
}

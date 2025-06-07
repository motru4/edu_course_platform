package app

import (
	"fmt"
	"net"
	"os"
	"time"

	"auth-service/internal/config"
	"auth-service/internal/database"
	"auth-service/internal/repositories"
	"auth-service/internal/security/jwt"
	"auth-service/internal/security/password"
	"auth-service/internal/services"
	grpcserver "auth-service/internal/transport/grpc"
	"auth-service/internal/transport/http/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type App struct {
	cfg        *config.Config
	httpServer *gin.Engine
	grpcServer *grpc.Server
}

// @title Auth Service API
// @version 1.0
// @description Сервис аутентификации и авторизации
// @host localhost:8090
// @BasePath /api/v1/auth
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func New() (*App, error) {
	a := new(App)
	var err error

	a.cfg, err = config.New()
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		MaxAge:           12 * time.Hour,
	}

	gin.SetMode(os.Getenv("GIN_MODE"))

	a.httpServer = gin.Default()

	a.httpServer.Use(
		cors.New(corsConfig),
		gin.Recovery(),
	)

	// Connect to database
	db, err := database.NewPostgresDB(a.cfg.Database.URL)
	if err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	refreshRepo := repositories.NewRefreshRepository(db)
	verificationRepo := repositories.NewVerificationRepository(db)

	// Initialize services
	emailService := services.NewEmailService(a.cfg.SMTP)
	tokenManager := jwt.NewJWTManager(a.cfg.Token)
	passwordHasher := password.NewHasher(a.cfg.Security.PasswordPepper)

	authService := services.NewAuthService(
		userRepo,
		refreshRepo,
		verificationRepo,
		emailService,
		tokenManager,
		passwordHasher,
		a.cfg,
	)

	// Initialize gRPC server
	a.grpcServer = grpcserver.NewServer(authService)

	// Initialize HTTP handlers
	handler.NewHandler(a.httpServer, authService, a.cfg)

	return a, nil
}

func (a *App) Run() error {
	// Start gRPC server
	lis, err := net.Listen("tcp", ":"+a.cfg.Server.GRPCPort)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	go func() {
		fmt.Printf("Starting gRPC server on port %s\n", a.cfg.Server.GRPCPort)
		if err := a.grpcServer.Serve(lis); err != nil {
			fmt.Printf("gRPC server error: %s\n", err)
		}
	}()

	// Start HTTP server
	fmt.Printf("Starting HTTP server on port %s\n", a.cfg.Server.Port)
	if err := a.httpServer.Run(":" + a.cfg.Server.Port); err != nil {
		return fmt.Errorf("http server error: %w", err)
	}

	return nil
}

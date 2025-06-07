package app

import (
	"context"
	"fmt"
	"os"
	"time"

	_ "game/docs" // Импорт для инициализации Swagger документации
	"game/internal/config"
	"game/internal/database"
	"game/internal/repositories"
	"game/internal/services"
	"game/internal/transport/http/handler"
	"game/internal/transport/http/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	cfg        *config.Config
	httpServer *gin.Engine
	grpcConn   *grpc.ClientConn
}

// @title Игровая платформа API
// @version 1.0
// @description API для игровой платформы с мини-играми и таблицами лидеров
// @description Включает функционал для игры-кликера, отслеживания статистики и просмотра таблицы лидеров
// @host localhost:8090
// @BasePath /api/v1/game
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @schemes http https
// @tags.name game "Игры"
// @tags.description game "Операции с мини-играми (кликер и др.)"
func New() (*App, error) {
	a := &App{}

	// Загрузка конфигурации
	cfg, err := config.New()
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки конфигурации: %w", err)
	}
	a.cfg = cfg

	// Подключение к gRPC серверу авторизации
	a.grpcConn, err = grpc.Dial(
		a.cfg.Auth.GRPCAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к серверу авторизации: %w", err)
	}

	// Настройка CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.MaxAge = 12 * time.Hour

	// Настройка Gin
	gin.SetMode(os.Getenv("GIN_MODE"))
	a.httpServer = gin.Default()
	a.httpServer.Use(
		cors.New(corsConfig),
		gin.Recovery(),
		gin.Logger(),
	)

	// Подключение к базе данных
	db, err := database.NewPostgresDB(a.cfg.Database.URL)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	// Инициализация репозиториев
	clickerRepo := repositories.NewClickerRepository(db)

	// Инициализация сервисов
	clickerService := services.NewClickerService(clickerRepo)

	// Инициализация HTTP обработчиков
	authRolesMiddleware := middleware.NewRolesMiddleware(a.grpcConn)
	handler.NewClickerHandler(a.httpServer, clickerService, authRolesMiddleware)

	// Настройка Swagger
	a.httpServer.GET("/api/v1/game/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return a, nil
}

// Run запускает приложение
func (a *App) Run() error {
	fmt.Printf("Запуск HTTP сервера на порту %s\n", a.cfg.Server.Port)
	if err := a.httpServer.Run(":" + a.cfg.Server.Port); err != nil {
		return fmt.Errorf("ошибка HTTP сервера: %w", err)
	}
	return nil
}

// Stop останавливает приложение
func (a *App) Stop(ctx context.Context) error {
	// Закрываем gRPC соединение
	if a.grpcConn != nil {
		if err := a.grpcConn.Close(); err != nil {
			return err
		}
	}

	return nil
}

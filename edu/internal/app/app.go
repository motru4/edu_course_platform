package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "course2/docs" // Импорт для инициализации Swagger документации
	"course2/internal/config"
	"course2/internal/database"
	"course2/internal/repositories"
	"course2/internal/services"
	"course2/internal/transport/http/handler"
	"course2/internal/transport/http/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	cfg        *config.Config
	router     *gin.Engine
	httpServer *http.Server
	grpcConn   *grpc.ClientConn
}

// @title Образовательная платформа API
// @version 1.0
// @description API для образовательной платформы с курсами, тестами и системой прогресса обучения
// @description Включает функционал для студентов (просмотр и покупка курсов, прохождение уроков и тестов),
// @description администраторов (управление курсами и модерация) и общедоступные эндпоинты (категории, публичная информация о курсах).
// @host localhost:8081
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @schemes http https
// @tags.name courses "Курсы"
// @tags.description courses "Операции с курсами (просмотр, создание, редактирование)"
// @tags.name lessons "Уроки"
// @tags.description lessons "Операции с уроками и тестами"
// @tags.name categories "Категории"
// @tags.description categories "Операции с категориями курсов"
// @tags.name admin "Администрирование"
// @tags.description admin "Административные операции (модерация, управление контентом)"
// @tags.name profile "Профиль"
// @tags.description profile "Операции с профилем пользователя"
// @tags.name progress "Прогресс"
// @tags.description progress "Отслеживание прогресса обучения"
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
	a.router = gin.Default()
	a.router.Use(
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
	courseRepo := repositories.NewCourseRepository(db)
	lessonRepo := repositories.NewLessonRepository(db)
	testRepo := repositories.NewTestRepository(db)
	reviewRepo := repositories.NewReviewRepository(db)
	progressRepo := repositories.NewProgressRepository(db)
	purchaseRepo := repositories.NewPurchaseRepository(db)
	userRepo := repositories.NewUserRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)

	// Инициализация сервисов
	courseService := services.NewCourseService(courseRepo, lessonRepo, testRepo, reviewRepo)
	userService := services.NewUserService(userRepo, purchaseRepo)
	paymentService := services.NewPaymentService(courseRepo, purchaseRepo)
	moderationService := services.NewModerationService(courseRepo)
	categoryService := services.NewCategoryService(categoryRepo)

	// Инициализация HTTP обработчиков
	authRolesMiddleware := middleware.NewRolesMiddleware(a.grpcConn)
	handler.NewCourseHandler(a.router, courseService, authRolesMiddleware)
	handler.NewStudentHandler(a.router, courseService, paymentService, progressRepo, purchaseRepo, authRolesMiddleware)
	handler.NewProgressHandler(a.router, courseService, progressRepo, authRolesMiddleware)
	handler.NewProfileHandler(a.router, userService, authRolesMiddleware)
	handler.NewAdminHandler(a.router, moderationService, authRolesMiddleware)
	handler.NewCategoryHandler(a.router, categoryService)

	// Настройка Swagger
	a.router.GET("/api/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Настройка HTTP сервера
	a.httpServer = &http.Server{
		Addr:    ":" + a.cfg.Server.Port,
		Handler: a.router,
	}

	return a, nil
}

// Run запускает приложение
func (a *App) Run() error {
	fmt.Printf("Запуск HTTP сервера на порту %s\n", a.cfg.Server.Port)
	if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("ошибка HTTP сервера: %w", err)
	}
	return nil
}

// Stop останавливает приложение
func (a *App) Stop(ctx context.Context) error {
	var err error

	// Закрываем HTTP сервер
	if a.httpServer != nil {
		if err = a.httpServer.Shutdown(ctx); err != nil {
			return err
		}
	}

	// Закрываем gRPC соединение
	if a.grpcConn != nil {
		if err = a.grpcConn.Close(); err != nil {
			return err
		}
	}

	return nil
}

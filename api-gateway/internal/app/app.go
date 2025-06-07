package app

import (
	"api-gateway/internal/config"
	"api-gateway/internal/proxy"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// App представляет собой API-шлюз
type App struct {
	cfg        *config.Config
	httpServer *gin.Engine
}

// New создает новый экземпляр App
func New() (*App, error) {
	app := new(App)

	// Загружаем конфигурацию
	app.cfg = config.New()

	// Настраиваем CORS
	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		MaxAge:           12 * time.Hour,
	}

	// Включаем режим отладки для Gin
	gin.SetMode(gin.DebugMode)

	// Создаем HTTP-сервер с Gin
	app.httpServer = gin.Default()

	// Добавляем middleware
	app.httpServer.Use(
		cors.New(corsConfig),
		gin.Recovery(),
		gin.Logger(),
	)

	// Добавляем простой хелсчек
	app.httpServer.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Создаем и настраиваем прокси для сервисов
	serviceProxy := proxy.NewServiceProxy(app.cfg)
	serviceProxy.SetupRoutes(app.httpServer)

	fmt.Printf("API Gateway настроен на порту %s\n", app.cfg.Server.Port)

	return app, nil
}

// Run запускает API-шлюз
func (a *App) Run() error {
	fmt.Printf("Запуск API-шлюза на порту %s\n", a.cfg.Server.Port)
	return a.httpServer.Run(":" + a.cfg.Server.Port)
}

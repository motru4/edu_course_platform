package proxy

import (
	"api-gateway/internal/config"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// ServiceProxy представляет собой прокси для сервисов
type ServiceProxy struct {
	cfg *config.Config
}

// NewServiceProxy создает новый экземпляр ServiceProxy
func NewServiceProxy(cfg *config.Config) *ServiceProxy {
	return &ServiceProxy{
		cfg: cfg,
	}
}

// ProxyRequest перенаправляет запрос к соответствующему сервису
func (p *ServiceProxy) ProxyRequest(c *gin.Context, serviceURL string) {
	// Проверяем и корректируем URL сервиса
	if !strings.HasPrefix(serviceURL, "http://") && !strings.HasPrefix(serviceURL, "https://") {
		serviceURL = "http://" + serviceURL
	}

	remote, err := url.Parse(serviceURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при разборе URL сервиса"})
		return
	}

	// Создаем реверс-прокси
	proxy := httputil.NewSingleHostReverseProxy(remote)

	// Настраиваем директор для модификации запроса
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		// Сохраняем оригинальный путь
		req.URL.Path = c.Request.URL.Path
		req.URL.RawQuery = c.Request.URL.RawQuery

		// Копируем заголовки из исходного запроса
		for k, v := range c.Request.Header {
			req.Header[k] = v
		}
	}

	// Настраиваем обработчик ошибок
	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "ошибка при обращении к сервису",
		})
	}

	// Убираем заголовок Content-Length, чтобы не конфликтовал при передаче
	c.Request.Header.Del("Content-Length")

	// Проксируем запрос к целевому сервису
	proxy.ServeHTTP(c.Writer, c.Request)
}

// SetupRoutes настраивает маршруты для API-шлюза
func (p *ServiceProxy) SetupRoutes(router *gin.Engine) {
	// Тестовый маршрут для проверки работоспособности
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Настраиваем маршруты для auth сервиса
	authGroup := router.Group("/api/v1/auth")
	authGroup.Any("/*path", func(c *gin.Context) {
		p.ProxyRequest(c, p.cfg.Services.AuthService.URL)
	})

	// Настраиваем маршруты для edu сервиса
	eduGroup := router.Group("/api/v1/edu")
	eduGroup.Any("/*path", func(c *gin.Context) {
		p.ProxyRequest(c, p.cfg.Services.EduService.URL)
	})

	// Настраиваем маршруты для game сервиса
	gameGroup := router.Group("/api/v1/game")
	gameGroup.Any("/*path", func(c *gin.Context) {
		p.ProxyRequest(c, p.cfg.Services.GameService.URL)
	})
}

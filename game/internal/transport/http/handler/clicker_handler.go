package handler

import (
	"log"
	"net/http"
	"strconv"

	"game/internal/models"
	"game/internal/services"
	"game/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ClickerHandler представляет обработчик для API кликера
type ClickerHandler struct {
	service *services.ClickerService
}

// NewClickerHandler создает новый экземпляр ClickerHandler и регистрирует маршруты
func NewClickerHandler(router *gin.Engine, service *services.ClickerService, authMiddleware *middleware.RolesMiddleware) {
	handler := &ClickerHandler{
		service: service,
	}

	gameGroup := router.Group("/api/v1/game")

	// Маршруты, требующие авторизации
	authorized := gameGroup.Group("/")
	authorized.Use(authMiddleware.RequireRoles("student", "admin"))
	{
		// Сохранение кликов
		authorized.POST("/clicker/clicks", handler.SaveClicks)

		// Получение статистики пользователя
		authorized.GET("/clicker/stats", handler.GetStats)

		// Получение таблицы лидеров
		authorized.GET("/clicker/leaderboard", handler.GetLeaderboard)
	}
}

// SaveClicks godoc
// @Summary Сохранить клики пользователя
// @Description Сохраняет количество кликов пользователя и обновляет его статистику
// @Tags game
// @Accept json
// @Produce json
// @Param request body models.ClickRequest true "Данные о кликах"
// @Success 200 {object} models.ClickResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /clicker/clicks [post]
// @Security BearerAuth
func (h *ClickerHandler) SaveClicks(c *gin.Context) {
	var req models.ClickRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверный формат запроса"})
		return
	}

	// Получаем ID пользователя из контекста
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "пользователь не аутентифицирован"})
		return
	}

	// Сохраняем клики
	response, err := h.service.SaveClicks(c.Request.Context(), userID.(uuid.UUID), &req)
	if err != nil {
		switch err {
		case services.ErrInvalidClickRate:
			c.JSON(http.StatusBadRequest, gin.H{"message": "обнаружена подозрительная активность: слишком высокая скорость кликов"})
		case services.ErrInvalidSession:
			c.JSON(http.StatusBadRequest, gin.H{"message": "обнаружена подозрительная активность: недопустимое время сессии"})
		case services.ErrInvalidTime:
			c.JSON(http.StatusBadRequest, gin.H{"message": "обнаружена подозрительная активность: несоответствие времени"})
		default:
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "ошибка при сохранении кликов"})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetStats godoc
// @Summary Получить статистику пользователя
// @Description Возвращает статистику пользователя в игре-кликере
// @Tags game
// @Produce json
// @Success 200 {object} models.StatsResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /clicker/stats [get]
// @Security BearerAuth
func (h *ClickerHandler) GetStats(c *gin.Context) {
	// Получаем ID пользователя из контекста
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "пользователь не аутентифицирован"})
		return
	}

	// Получаем статистику
	stats, err := h.service.GetStats(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "ошибка при получении статистики"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetLeaderboard godoc
// @Summary Получить таблицу лидеров
// @Description Возвращает таблицу лидеров игры-кликера
// @Tags game
// @Produce json
// @Param limit query int false "Лимит записей (по умолчанию 10)"
// @Param offset query int false "Смещение (по умолчанию 0)"
// @Success 200 {object} models.LeaderboardResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /clicker/leaderboard [get]
// @Security BearerAuth
func (h *ClickerHandler) GetLeaderboard(c *gin.Context) {
	// Получаем ID пользователя из контекста
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "пользователь не аутентифицирован"})
		return
	}

	// Получаем параметры запроса
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Получаем таблицу лидеров
	leaderboard, err := h.service.GetLeaderboard(c.Request.Context(), userID.(uuid.UUID), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "ошибка при получении таблицы лидеров"})
		return
	}

	c.JSON(http.StatusOK, leaderboard)
}

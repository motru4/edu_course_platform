package handler

import (
	"course2/internal/models"
	"course2/internal/services"
	"course2/internal/transport/http/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProfileHandler struct {
	userService *services.UserService
}

func NewProfileHandler(router *gin.Engine, userService *services.UserService, authMiddleware *middleware.RolesMiddleware) {
	handler := &ProfileHandler{
		userService: userService,
	}

	profile := router.Group("/api/v1/edu/profile")
	profile.Use(authMiddleware.RequireRoles("student", "admin"))
	{
		profile.GET("", handler.GetProfile)
		profile.PUT("", handler.UpdateProfile)
		profile.GET("/courses", handler.GetPurchasedCourses)
		profile.GET("/xp", handler.GetTotalXP)
	}
}

// @Summary Получить профиль
// @Description Получить профиль текущего пользователя
// @Tags profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserProfile
// @Router /profile [get]
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	profile, err := h.userService.GetProfile(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// @Summary Обновить профиль
// @Description Обновить профиль текущего пользователя
// @Tags profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profile body models.UserProfile true "Данные профиля"
// @Success 200 {object} models.UserProfile
// @Router /profile [put]
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	var profile models.UserProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile.ID = userID.(uuid.UUID)
	if err := h.userService.UpdateProfile(c.Request.Context(), &profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// @Summary Получить купленные курсы
// @Description Получить список купленных курсов текущего пользователя
// @Tags profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Course
// @Router /profile/courses [get]
func (h *ProfileHandler) GetPurchasedCourses(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	courses, err := h.userService.GetPurchasedCourses(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, courses)
}

// @Summary Получить общий XP
// @Description Получить общее количество XP текущего пользователя
// @Tags profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]int
// @Router /profile/xp [get]
func (h *ProfileHandler) GetTotalXP(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	xp, err := h.userService.GetTotalXP(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_xp": xp})
}

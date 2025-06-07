package handler

import (
	"course2/internal/models"
	"course2/internal/services"
	"course2/internal/transport/http/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminHandler struct {
	moderationService *services.ModerationService
}

func NewAdminHandler(router *gin.Engine, moderationService *services.ModerationService, authMiddleware *middleware.RolesMiddleware) {
	handler := &AdminHandler{
		moderationService: moderationService,
	}

	admin := router.Group("/api/v1/edu/admin")
	admin.Use(authMiddleware.RequireRoles("admin"))
	{
		admin.GET("/courses/pending", handler.ListPendingCourses)
		admin.POST("/courses/:id/approve", handler.ApproveCourse)
		admin.POST("/courses/:id/reject", handler.RejectCourse)

		// Управление курсами
		admin.POST("/courses", handler.CreateCourse)
		admin.PUT("/courses/:id", handler.UpdateCourse)
		admin.DELETE("/courses/:id", handler.DeleteCourse)
	}
}

// @Summary Список ожидающих модерации курсов
// @Description Получить список курсов, ожидающих модерации
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Номер страницы"
// @Param limit query int false "Количество элементов на странице"
// @Success 200 {array} models.Course
// @Router /admin/courses/pending [get]
func (h *AdminHandler) ListPendingCourses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	courses, err := h.moderationService.ListPendingCourses(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, courses)
}

// @Summary Одобрить курс
// @Description Одобрить курс для публикации
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID курса"
// @Success 200 {object} map[string]string
// @Router /admin/courses/{id}/approve [post]
func (h *AdminHandler) ApproveCourse(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID"})
		return
	}

	if err := h.moderationService.ApproveCourse(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Курс одобрен"})
}

// @Summary Отклонить курс
// @Description Отклонить курс с указанием причины
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID курса"
// @Param reason body map[string]string true "Причина отклонения"
// @Success 200 {object} map[string]string
// @Router /admin/courses/{id}/reject [post]
func (h *AdminHandler) RejectCourse(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID"})
		return
	}

	var body struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.moderationService.RejectCourse(c.Request.Context(), id, body.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Курс отклонен"})
}

// @Summary Создать курс
// @Description Создать новый курс
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param course body models.Course true "Данные курса"
// @Success 201 {object} models.Course
// @Router /admin/courses [post]
func (h *AdminHandler) CreateCourse(c *gin.Context) {
	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.moderationService.CreateCourse(c.Request.Context(), &course); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, course)
}

// @Summary Обновить курс
// @Description Обновить существующий курс
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID курса"
// @Param course body models.Course true "Данные курса"
// @Success 200 {object} models.Course
// @Router /admin/courses/{id} [put]
func (h *AdminHandler) UpdateCourse(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID"})
		return
	}

	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course.ID = id
	if err := h.moderationService.UpdateCourse(c.Request.Context(), &course); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, course)
}

// @Summary Удалить курс
// @Description Удалить существующий курс
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID курса"
// @Success 204 "No Content"
// @Router /admin/courses/{id} [delete]
func (h *AdminHandler) DeleteCourse(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID"})
		return
	}

	if err := h.moderationService.DeleteCourse(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

package handler

import (
	"course2/internal/services"
	"course2/internal/transport/http/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CourseHandler struct {
	courseService *services.CourseService
}

func NewCourseHandler(router *gin.Engine, courseService *services.CourseService, authMiddleware *middleware.RolesMiddleware) {
	handler := &CourseHandler{
		courseService: courseService,
	}

	courses := router.Group("/api/v1/edu/courses")
	{
		// Публичные эндпоинты
		courses.GET("", handler.ListCourses)
		courses.GET("/:id", handler.GetCourse)
		courses.GET("/category/:categoryId", handler.ListCoursesByCategory)

		// Защищенные эндпоинты
		/* authorized := courses.Group("")
		authorized.Use(authMiddleware.RequireRoles("admin"))
		{
			authorized.POST("", handler.CreateCourse)
			authorized.PUT("/:id", handler.UpdateCourse)
			authorized.DELETE("/:id", handler.DeleteCourse)
		} */
	}
}

// @Summary Список курсов
// @Description Получить список всех курсов
// @Tags courses
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы"
// @Param limit query int false "Количество элементов на странице"
// @Success 200 {array} models.Course
// @Router /courses [get]
func (h *CourseHandler) ListCourses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	courses, err := h.courseService.ListCourses(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, courses)
}

// @Summary Получить курс
// @Description Получить информацию о курсе по ID
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "ID курса"
// @Success 200 {object} models.Course
// @Router /courses/{id} [get]
func (h *CourseHandler) GetCourse(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID"})
		return
	}

	course, err := h.courseService.GetCourse(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if course == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Курс не найден"})
		return
	}

	c.JSON(http.StatusOK, course)
}

// @Summary Список курсов по категории
// @Description Получить список курсов в указанной категории
// @Tags courses
// @Accept json
// @Produce json
// @Param categoryId path string true "ID категории"
// @Param page query int false "Номер страницы"
// @Param limit query int false "Количество элементов на странице"
// @Success 200 {array} models.Course
// @Router /courses/category/{categoryId} [get]
func (h *CourseHandler) ListCoursesByCategory(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("categoryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID категории"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	courses, err := h.courseService.ListCoursesByCategory(c.Request.Context(), categoryID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, courses)
}

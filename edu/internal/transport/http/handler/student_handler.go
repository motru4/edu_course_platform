package handler

import (
	"course2/internal/models"
	"course2/internal/repositories"
	"course2/internal/services"
	"course2/internal/transport/http/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StudentHandler struct {
	courseService  *services.CourseService
	paymentService *services.PaymentService
	progressRepo   *repositories.ProgressRepository
	purchaseRepo   *repositories.PurchaseRepository
}

func NewStudentHandler(
	router *gin.Engine,
	courseService *services.CourseService,
	paymentService *services.PaymentService,
	progressRepo *repositories.ProgressRepository,
	purchaseRepo *repositories.PurchaseRepository,
	authMiddleware *middleware.RolesMiddleware,
) {
	handler := &StudentHandler{
		courseService:  courseService,
		paymentService: paymentService,
		progressRepo:   progressRepo,
		purchaseRepo:   purchaseRepo,
	}

	student := router.Group("/api/v1/student")
	{
		// Публичные эндпоинты для курсов
		student.GET("/courses/:courseId/lessons", handler.GetCourseLessons)
		student.GET("/lessons/:lessonId", handler.GetLesson)

		// Защищенные эндпоинты
		authorized := student.Group("")
		authorized.Use(authMiddleware.RequireRoles("student"))
		{
			authorized.POST("/courses/purchase", handler.PurchaseCourse)
			authorized.GET("/courses/:courseId/structure", handler.GetCourseStructure)
			authorized.GET("/lessons/:lessonId/test", handler.GetLessonTest)
		}
	}
}

// @Summary Получить уроки курса
// @Description Получить список уроков для конкретного курса
// @Tags courses
// @Accept json
// @Produce json
// @Param courseId path string true "ID курса"
// @Success 200 {array} models.Lesson
// @Router /student/courses/{courseId}/lessons [get]
func (h *StudentHandler) GetCourseLessons(c *gin.Context) {
	courseID, err := uuid.Parse(c.Param("courseId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID курса"})
		return
	}

	// Проверяем существование курса
	course, err := h.courseService.GetCourse(c.Request.Context(), courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if course == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Курс не найден"})
		return
	}

	lessons, err := h.courseService.GetLessons(c.Request.Context(), courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Если пользователь авторизован, добавляем информацию о прогрессе
	if userID, exists := c.Get("user_id"); exists {
		for _, lesson := range lessons {
			progress, err := h.progressRepo.GetLessonProgress(c.Request.Context(), userID.(uuid.UUID), lesson.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if progress != nil {
				lesson.Completed = progress.ViewedAt != nil && (!lesson.RequiresTest || progress.PassedTest)
			}
		}
	}

	c.JSON(http.StatusOK, lessons)
}

// @Summary Получить урок
// @Description Получить содержимое конкретного урока
// @Tags lessons
// @Accept json
// @Produce json
// @Param lessonId path string true "ID урока"
// @Success 200 {object} models.Lesson
// @Router /student/lessons/{lessonId} [get]
func (h *StudentHandler) GetLesson(c *gin.Context) {
	lessonID, err := uuid.Parse(c.Param("lessonId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID урока"})
		return
	}

	lesson, err := h.courseService.GetLesson(c.Request.Context(), lessonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if lesson == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Урок не найден"})
		return
	}

	// Проверяем доступ к уроку
	if userID, exists := c.Get("user_id"); exists {
		purchased, err := h.purchaseRepo.HasPurchased(c.Request.Context(), userID.(uuid.UUID), lesson.CourseID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !purchased {
			c.JSON(http.StatusForbidden, gin.H{"error": "Необходимо купить курс для доступа к уроку"})
			return
		}
	}

	c.JSON(http.StatusOK, lesson)
}

// PurchaseCourseRequest модель запроса для покупки курса
type PurchaseCourseRequest struct {
	CourseID uuid.UUID `json:"course_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// @Summary Купить курс
// @Description Покупка курса пользователем
// @Tags courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body PurchaseCourseRequest true "Данные для покупки курса"
// @Success 200 {object} map[string]interface{}
// @Router /student/courses/purchase [post]
func (h *StudentHandler) PurchaseCourse(c *gin.Context) {
	var request PurchaseCourseRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)

	// Проверяем, не куплен ли уже курс
	purchased, err := h.purchaseRepo.HasPurchased(c.Request.Context(), userID, request.CourseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if purchased {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Курс уже куплен"})
		return
	}

	// Обрабатываем покупку через платежный сервис
	if err := h.paymentService.PurchaseCourse(c.Request.Context(), userID, request.CourseID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Курс успешно куплен",
	})
}

// @Summary Получить структуру курса
// @Description Получить полную структуру курса с уроками и прогрессом для купленного курса
// @Tags courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param courseId path string true "ID курса"
// @Success 200 {object} models.CourseStructure
// @Router /student/courses/{courseId}/structure [get]
func (h *StudentHandler) GetCourseStructure(c *gin.Context) {
	courseID, err := uuid.Parse(c.Param("courseId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID курса"})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)

	// Проверяем, куплен ли курс
	purchased, err := h.purchaseRepo.HasPurchased(c.Request.Context(), userID, courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !purchased {
		c.JSON(http.StatusForbidden, gin.H{"error": "Необходимо купить курс для доступа к структуре"})
		return
	}

	// Получаем курс
	course, err := h.courseService.GetCourse(c.Request.Context(), courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if course == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Курс не найден"})
		return
	}

	// Получаем уроки
	lessons, err := h.courseService.GetLessons(c.Request.Context(), courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Получаем прогресс по всем урокам
	for _, lesson := range lessons {
		progress, err := h.progressRepo.GetLessonProgress(c.Request.Context(), userID, lesson.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if progress != nil {
			lesson.Completed = progress.ViewedAt != nil && (!lesson.RequiresTest || progress.PassedTest)
			lesson.TestScore = progress.TestScore
			lesson.PassedTest = progress.PassedTest
			lesson.ViewedAt = progress.ViewedAt
		}

		// Если урок требует тест, добавляем информацию о его наличии
		if lesson.RequiresTest {
			test, err := h.courseService.GetTestByLessonID(c.Request.Context(), lesson.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if test != nil {
				lesson.HasTest = true
			}
		}
	}

	// Получаем общий прогресс по курсу
	totalLessons := len(lessons)
	completedLessons := 0
	for _, lesson := range lessons {
		if lesson.Completed {
			completedLessons++
		}
	}

	response := &models.CourseStructure{
		Course:           course,
		Lessons:          lessons,
		TotalLessons:     totalLessons,
		CompletedLessons: completedLessons,
	}
	response.CalculateProgress()

	c.JSON(http.StatusOK, response)
}

// @Summary Получить тест урока
// @Description Получить тест, привязанный к уроку (если есть)
// @Tags lessons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param lessonId path string true "ID урока"
// @Success 200 {object} models.TestResponse
// @Router /student/lessons/{lessonId}/test [get]
func (h *StudentHandler) GetLessonTest(c *gin.Context) {
	lessonID, err := uuid.Parse(c.Param("lessonId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID урока"})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)

	// Получаем урок для проверки принадлежности к курсу
	lesson, err := h.courseService.GetLesson(c.Request.Context(), lessonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if lesson == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Урок не найден"})
		return
	}

	// Проверяем, куплен ли курс
	purchased, err := h.purchaseRepo.HasPurchased(c.Request.Context(), userID, lesson.CourseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !purchased {
		c.JSON(http.StatusForbidden, gin.H{"error": "Необходимо купить курс для доступа к тесту"})
		return
	}

	// Получаем тест
	test, err := h.courseService.GetTestByLessonID(c.Request.Context(), lessonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if test == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Тест не найден"})
		return
	}

	// Получаем вопросы теста
	questions, err := h.courseService.GetTestQuestions(c.Request.Context(), test.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Получаем прогресс по уроку
	progress, err := h.progressRepo.GetLessonProgress(c.Request.Context(), userID, lessonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Подготавливаем вопросы для отправки (удаляем правильные ответы)
	testQuestions := make([]*models.Question, len(questions))
	for i, q := range questions {
		testQuestions[i] = &models.Question{
			ID:           q.ID,
			QuestionText: q.QuestionText,
			Options:      q.Options,
		}
	}

	response := &models.TestResponse{
		Test:          test,
		Questions:     testQuestions,
		PassingScore:  test.PassingScore,
		LastScore:     nil,
		Passed:        false,
		AttemptsCount: 0,
	}

	// Если есть прогресс, обновляем значения
	if progress != nil {
		response.LastScore = progress.TestScore
		response.Passed = progress.PassedTest
		response.AttemptsCount = progress.AttemptsCount
	}

	c.JSON(http.StatusOK, response)
}

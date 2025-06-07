package handler

import (
	"course2/internal/models"
	"course2/internal/repositories"
	"course2/internal/services"
	"course2/internal/transport/http/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProgressHandler struct {
	courseService *services.CourseService
	progressRepo  *repositories.ProgressRepository
}

func NewProgressHandler(
	router *gin.Engine,
	courseService *services.CourseService,
	progressRepo *repositories.ProgressRepository,
	authMiddleware *middleware.RolesMiddleware,
) {
	handler := &ProgressHandler{
		courseService: courseService,
		progressRepo:  progressRepo,
	}

	progress := router.Group("/api/v1/edu/progress")
	progress.Use(authMiddleware.RequireRoles("student"))
	{
		progress.GET("/courses/:courseId", handler.GetCourseProgress)
		progress.POST("/lessons/:lessonId/view", handler.MarkLessonViewed)
		progress.POST("/lessons/:lessonId/test", handler.SubmitTest)
	}
}

// @Summary Получить прогресс по курсу
// @Description Получить прогресс пользователя по конкретному курсу
// @Tags progress
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param courseId path string true "ID курса"
// @Success 200 {object} models.CourseProgress
// @Router /progress/courses/{courseId} [get]
func (h *ProgressHandler) GetCourseProgress(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	courseID, err := uuid.Parse(c.Param("courseId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID курса"})
		return
	}

	progress, err := h.progressRepo.GetCourseProgress(c.Request.Context(), userID.(uuid.UUID), courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// @Summary Отметить урок как просмотренный
// @Description Отметить урок как просмотренный и начислить XP
// @Tags progress
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param lessonId path string true "ID урока"
// @Success 200 {object} map[string]interface{}
// @Router /progress/lessons/{lessonId}/view [post]
func (h *ProgressHandler) MarkLessonViewed(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	lessonID, err := uuid.Parse(c.Param("lessonId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID урока"})
		return
	}

	// Получаем информацию об уроке
	lesson, err := h.courseService.GetLesson(c.Request.Context(), lessonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Проверяем существующий прогресс
	existingProgress, err := h.progressRepo.GetLessonProgress(c.Request.Context(), userID.(uuid.UUID), lessonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()

	if existingProgress != nil {
		// Обновляем существующую запись
		if existingProgress.ViewedAt == nil {
			existingProgress.ViewedAt = &now
		}
		if !lesson.RequiresTest && existingProgress.CompletedAt == nil {
			existingProgress.CompletedAt = &now
			existingProgress.IsCompleted = true
		}

		if err := h.progressRepo.UpdateLessonProgress(c.Request.Context(), existingProgress); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// Создаем новую запись
		progress := &models.LessonProgress{
			ID:          uuid.New(),
			UserID:      userID.(uuid.UUID),
			LessonID:    lessonID,
			ViewedAt:    &now,
			IsCompleted: !lesson.RequiresTest,
		}

		// Если урок не требует теста, устанавливаем completed_at
		if !lesson.RequiresTest {
			progress.CompletedAt = &now
		}

		if err := h.progressRepo.CreateLessonProgress(c.Request.Context(), progress); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Начисляем XP за просмотр урока только при первом просмотре
		xpEntry := &models.XPEntry{
			ID:       uuid.New(),
			UserID:   userID.(uuid.UUID),
			CourseID: lesson.CourseID,
			LessonID: lessonID,
			Type:     "lesson_view",
			Amount:   5,
		}

		if err := h.progressRepo.AddXPEntry(c.Request.Context(), xpEntry); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	xpAwarded := 0
	if existingProgress == nil {
		xpAwarded = 5
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Урок отмечен как просмотренный",
		"xp_awarded": xpAwarded,
	})
}

// @Summary Отправить ответы на тест
// @Description Отправить ответы на тест и получить результат
// @Tags progress
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param lessonId path string true "ID урока"
// @Param answers body map[string]int true "Ответы на вопросы"
// @Success 200 {object} map[string]interface{}
// @Router /progress/lessons/{lessonId}/test [post]
func (h *ProgressHandler) SubmitTest(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	lessonID, err := uuid.Parse(c.Param("lessonId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID урока"})
		return
	}

	var answers map[string]int
	if err := c.ShouldBindJSON(&answers); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем тест для урока
	test, err := h.courseService.GetTestByLessonID(c.Request.Context(), lessonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if test == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Тест не найден"})
		return
	}

	// Получаем информацию об уроке для course_id
	lesson, err := h.courseService.GetLesson(c.Request.Context(), lessonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Проверяем ответы и вычисляем результат
	questions, err := h.courseService.GetTestQuestions(c.Request.Context(), test.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	correctAnswers := 0
	for _, q := range questions {
		if answer, ok := answers[q.ID.String()]; ok && answer == q.CorrectAnswer {
			correctAnswers++
		}
	}

	score := (correctAnswers * 100) / len(questions)
	passed := score >= test.PassingScore

	// Получаем текущий прогресс
	currentProgress, err := h.progressRepo.GetLessonProgress(c.Request.Context(), userID.(uuid.UUID), lessonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	progress := &models.LessonProgress{
		ID:            uuid.New(),
		UserID:        userID.(uuid.UUID),
		LessonID:      lessonID,
		TestScore:     &score,
		PassedTest:    passed,
		LastAttemptAt: &now,
		IsCompleted:   passed,
		AttemptsCount: 1,
	}

	if currentProgress != nil {
		progress.ID = currentProgress.ID
		progress.AttemptsCount = currentProgress.AttemptsCount + 1
		progress.ViewedAt = currentProgress.ViewedAt
	}

	// Если тест пройден, устанавливаем completed_at
	if passed {
		progress.CompletedAt = &now
	}

	if err := h.progressRepo.UpdateLessonProgress(c.Request.Context(), progress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Если тест пройден, начисляем XP
	if passed {
		xpEntry := &models.XPEntry{
			ID:       uuid.New(),
			UserID:   userID.(uuid.UUID),
			CourseID: lesson.CourseID,
			LessonID: lessonID,
			Type:     "test_pass",
			Amount:   10,
		}

		if err := h.progressRepo.AddXPEntry(c.Request.Context(), xpEntry); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":       "Тест успешно пройден",
			"score":         score,
			"xp_awarded":    10,
			"attempts_used": progress.AttemptsCount,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message":       "Тест не пройден",
			"score":         score,
			"attempts_used": progress.AttemptsCount,
		})
	}
}

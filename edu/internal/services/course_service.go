package services

import (
	"context"
	"course2/internal/models"
	"course2/internal/repositories"

	"github.com/google/uuid"
)

type CourseService struct {
	courseRepo *repositories.CourseRepository
	lessonRepo *repositories.LessonRepository
	testRepo   *repositories.TestRepository
	reviewRepo *repositories.ReviewRepository
}

func NewCourseService(
	courseRepo *repositories.CourseRepository,
	lessonRepo *repositories.LessonRepository,
	testRepo *repositories.TestRepository,
	reviewRepo *repositories.ReviewRepository,
) *CourseService {
	return &CourseService{
		courseRepo: courseRepo,
		lessonRepo: lessonRepo,
		testRepo:   testRepo,
		reviewRepo: reviewRepo,
	}
}

func (s *CourseService) GetCourse(ctx context.Context, id uuid.UUID) (*models.Course, error) {
	return s.courseRepo.GetByID(ctx, id)
}

func (s *CourseService) ListCourses(ctx context.Context, page, limit int) ([]*models.Course, error) {
	offset := (page - 1) * limit
	return s.courseRepo.List(ctx, offset, limit)
}

func (s *CourseService) ListCoursesByCategory(ctx context.Context, categoryID uuid.UUID, page, limit int) ([]*models.Course, error) {
	offset := (page - 1) * limit
	return s.courseRepo.ListByCategory(ctx, categoryID, offset, limit)
}

// Дополнительные методы для работы с уроками
func (s *CourseService) AddLesson(ctx context.Context, lesson *models.Lesson) error {
	lesson.ID = uuid.New()
	return s.lessonRepo.Create(ctx, lesson)
}

func (s *CourseService) UpdateLesson(ctx context.Context, lesson *models.Lesson) error {
	return s.lessonRepo.Update(ctx, lesson)
}

func (s *CourseService) DeleteLesson(ctx context.Context, id uuid.UUID) error {
	return s.lessonRepo.Delete(ctx, id)
}

// Методы для работы с тестами
func (s *CourseService) AddTest(ctx context.Context, test *models.Test) error {
	test.ID = uuid.New()
	return s.testRepo.Create(ctx, test)
}

func (s *CourseService) UpdateTest(ctx context.Context, test *models.Test) error {
	return s.testRepo.Update(ctx, test)
}

func (s *CourseService) DeleteTest(ctx context.Context, id uuid.UUID) error {
	return s.testRepo.Delete(ctx, id)
}

// Методы для работы с вопросами теста
func (s *CourseService) AddQuestion(ctx context.Context, question *models.Question) error {
	question.ID = uuid.New()
	return s.testRepo.CreateQuestion(ctx, question)
}

func (s *CourseService) UpdateQuestion(ctx context.Context, question *models.Question) error {
	return s.testRepo.UpdateQuestion(ctx, question)
}

func (s *CourseService) DeleteQuestion(ctx context.Context, id uuid.UUID) error {
	return s.testRepo.DeleteQuestion(ctx, id)
}

func (s *CourseService) GetTestByLessonID(ctx context.Context, lessonID uuid.UUID) (*models.Test, error) {
	return s.testRepo.GetByLessonID(ctx, lessonID)
}

func (s *CourseService) GetTestQuestions(ctx context.Context, testID uuid.UUID) ([]*models.Question, error) {
	return s.testRepo.GetQuestions(ctx, testID)
}

func (s *CourseService) GetLessons(ctx context.Context, courseID uuid.UUID) ([]*models.Lesson, error) {
	return s.lessonRepo.ListByCourse(ctx, courseID)
}

func (s *CourseService) GetLesson(ctx context.Context, lessonID uuid.UUID) (*models.Lesson, error) {
	return s.lessonRepo.GetByID(ctx, lessonID)
}

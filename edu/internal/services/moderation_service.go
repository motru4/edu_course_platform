package services

import (
	"context"
	"course2/internal/models"
	"course2/internal/repositories"

	"github.com/google/uuid"
)

type ModerationService struct {
	courseRepo *repositories.CourseRepository
}

func NewModerationService(courseRepo *repositories.CourseRepository) *ModerationService {
	return &ModerationService{
		courseRepo: courseRepo,
	}
}

func (s *ModerationService) ApproveCourse(ctx context.Context, courseID uuid.UUID) error {
	course, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return err
	}
	if course == nil {
		return ErrCourseNotFound
	}

	course.Status = "published"
	return s.courseRepo.Update(ctx, course)
}

func (s *ModerationService) RejectCourse(ctx context.Context, courseID uuid.UUID, reason string) error {
	course, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return err
	}
	if course == nil {
		return ErrCourseNotFound
	}

	course.Status = "rejected"
	return s.courseRepo.Update(ctx, course)
}

func (s *ModerationService) ListPendingCourses(ctx context.Context, page, limit int) ([]*models.Course, error) {
	offset := (page - 1) * limit
	return s.courseRepo.List(ctx, offset, limit)
}

// CreateCourse создает новый курс
func (s *ModerationService) CreateCourse(ctx context.Context, course *models.Course) error {
	course.ID = uuid.New()
	course.Status = "pending" // Новые курсы создаются со статусом "pending"
	return s.courseRepo.Create(ctx, course)
}

// UpdateCourse обновляет существующий курс
func (s *ModerationService) UpdateCourse(ctx context.Context, course *models.Course) error {
	existing, err := s.courseRepo.GetByID(ctx, course.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrCourseNotFound
	}

	// Сохраняем текущий статус курса
	course.Status = existing.Status
	return s.courseRepo.Update(ctx, course)
}

// DeleteCourse удаляет курс
func (s *ModerationService) DeleteCourse(ctx context.Context, courseID uuid.UUID) error {
	existing, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrCourseNotFound
	}

	return s.courseRepo.Delete(ctx, courseID)
}

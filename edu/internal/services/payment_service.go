package services

import (
	"context"
	"course2/internal/repositories"

	"github.com/google/uuid"
)

type PaymentService struct {
	courseRepo   *repositories.CourseRepository
	purchaseRepo *repositories.PurchaseRepository
}

func NewPaymentService(courseRepo *repositories.CourseRepository, purchaseRepo *repositories.PurchaseRepository) *PaymentService {
	return &PaymentService{
		courseRepo:   courseRepo,
		purchaseRepo: purchaseRepo,
	}
}

func (s *PaymentService) PurchaseCourse(ctx context.Context, userID, courseID uuid.UUID) error {
	// Проверяем, что курс существует
	course, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return err
	}
	if course == nil {
		return ErrCourseNotFound
	}

	// Проверяем, не куплен ли уже курс
	purchased, err := s.purchaseRepo.HasPurchased(ctx, userID, courseID)
	if err != nil {
		return err
	}
	if purchased {
		return ErrCourseAlreadyPurchased
	}

	// Здесь должна быть интеграция с платежной системой
	// ...

	// Записываем информацию о покупке
	return s.purchaseRepo.CreatePurchase(ctx, userID, courseID)
}

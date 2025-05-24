package services

import (
	"context"
	"course2/internal/models"
	"course2/internal/repositories"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo     *repositories.UserRepository
	purchaseRepo *repositories.PurchaseRepository
}

func NewUserService(userRepo *repositories.UserRepository, purchaseRepo *repositories.PurchaseRepository) *UserService {
	return &UserService{
		userRepo:     userRepo,
		purchaseRepo: purchaseRepo,
	}
}

func (s *UserService) GetProfile(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	totalXP, err := s.userRepo.GetTotalXP(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &models.UserProfile{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
		Role:      user.Role,
		TotalXP:   totalXP,
		Settings:  user.Settings,
	}, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, profile *models.UserProfile) error {
	return s.userRepo.UpdateProfile(ctx, profile)
}

func (s *UserService) GetPurchasedCourses(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	return s.purchaseRepo.GetPurchasedCourses(ctx, userID)
}

func (s *UserService) GetTotalXP(ctx context.Context, userID uuid.UUID) (int, error) {
	return s.userRepo.GetTotalXP(ctx, userID)
}

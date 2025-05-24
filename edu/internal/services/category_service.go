package services

import (
	"context"
	"course2/internal/models"
	"course2/internal/repositories"
)

type CategoryService struct {
	categoryRepo *repositories.CategoryRepository
}

func NewCategoryService(categoryRepo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *CategoryService) ListCategories(ctx context.Context) ([]*models.Category, error) {
	return s.categoryRepo.List(ctx)
}

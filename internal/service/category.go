package service

import (
	"CookFinder.Backend/internal/model"
	repository "CookFinder.Backend/internal/repo"
	"CookFinder.Backend/pkg/uuid"
	"context"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) Create(ctx context.Context, category *model.Category) error {
	if category.ID == "" {
		category.ID = uuid.V7().String()
	}

	return s.repo.Create(ctx, category)
}

func (s *CategoryService) GetAll(ctx context.Context) ([]model.Category, error) {
	return s.repo.GetAll(ctx)
}

func (s *CategoryService) GetByID(ctx context.Context, id string) (*model.Category, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CategoryService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *CategoryService) Update(ctx context.Context, category *model.Category) error {
	return s.repo.Update(ctx, category)
}

package service

import (
	"CookFinder.Backend/internal/model"
	repository "CookFinder.Backend/internal/repo"
	"CookFinder.Backend/pkg/uuid"
	"context"
)

type IngredientService struct {
	repo *repository.IngredientRepository
}

func NewIngredientService(repo *repository.IngredientRepository) *IngredientService {
	return &IngredientService{repo: repo}
}

func (s *IngredientService) Create(ctx context.Context, ingredient *model.Ingredient) error {
	if ingredient.ID == "" {
		ingredient.ID = uuid.V7().String()
	}

	return s.repo.Create(ctx, ingredient)
}

func (s *IngredientService) GetByID(ctx context.Context, id string) (*model.Ingredient, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *IngredientService) GetAll(ctx context.Context) ([]model.Ingredient, error) {
	return s.repo.GetAll(ctx)
}

func (s *IngredientService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

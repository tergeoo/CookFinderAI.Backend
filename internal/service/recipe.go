package service

import (
	"CookFinder.Backend/internal/model"
	repository "CookFinder.Backend/internal/repo"
	"CookFinder.Backend/pkg/uuid"
	"context"
)

type RecipeService struct {
	repo    *repository.RecipeRepository
	ingRepo *repository.IngredientRepository
	catRepo *repository.CategoryRepository
}

func NewRecipeService(
	repo *repository.RecipeRepository,
	ingRepo *repository.IngredientRepository,
	catRepo *repository.CategoryRepository,
) *RecipeService {
	return &RecipeService{repo: repo, ingRepo: ingRepo, catRepo: catRepo}
}

func (it *RecipeService) Create(ctx context.Context, recipe *model.Recipe) error {
	if recipe.ID == "" {
		recipe.ID = uuid.V7().String()
	}

	return it.repo.Create(ctx, recipe)
}

func (it *RecipeService) GetByID(ctx context.Context, id string) (*model.RecipeCategoryIngredients, error) {
	return it.repo.GetByID(ctx, id)
}

func (it *RecipeService) GetAll(ctx context.Context) ([]model.RecipeCategoryIngredients, error) {
	return it.repo.GetAll(ctx)
}

func (it *RecipeService) Update(ctx context.Context, recipe *model.Recipe) error {
	return it.repo.Update(ctx, recipe)
}

func (it *RecipeService) Delete(ctx context.Context, id string) error {
	return it.repo.Delete(ctx, id)
}

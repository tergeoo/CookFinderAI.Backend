package service

import (
	"CookFinder.Backend/internal/model"
	"CookFinder.Backend/internal/repo"
	"CookFinder.Backend/pkg/uuid"
	"context"
	"time"
)

type RecipeService struct {
	recipeRepo     *repo.RecipeRepository
	recipeIngrRepo *repo.RecipeIngredientRepository
}

func NewRecipeService(
	repo *repo.RecipeRepository,
	ingrRepo *repo.RecipeIngredientRepository,
) *RecipeService {
	return &RecipeService{
		recipeRepo:     repo,
		recipeIngrRepo: ingrRepo,
	}
}

func (s *RecipeService) CreateWithIngredients(ctx context.Context, recipe *model.Recipe, ingredients []model.RecipeIngredient) error {
	if recipe.ID == "" {
		recipe.ID = uuid.V7().String()
	}
	recipe.CreatedAt = time.Now()

	tx, err := s.recipeRepo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Сохраняем сам рецепт
	if err := s.recipeRepo.CreateWithTx(ctx, tx, recipe); err != nil {
		return err
	}

	// Сохраняем ингредиенты
	for _, ing := range ingredients {
		ing.RecipeID = recipe.ID
		if err := s.recipeIngrRepo.AddWithTx(ctx, tx, &ing); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *RecipeService) GetByID(ctx context.Context, id string) (*model.RecipeCategoryIngredients, error) {
	return s.recipeRepo.GetByID(ctx, id)
}

func (s *RecipeService) GetAll(ctx context.Context, search, categoryID string) ([]model.RecipeCategoryIngredients, error) {
	return s.recipeRepo.GetAll(ctx, search, categoryID)
}

func (s *RecipeService) Update(ctx context.Context, recipe *model.Recipe) error {
	return s.recipeRepo.Update(ctx, recipe)
}

func (s *RecipeService) Delete(ctx context.Context, id string) error {
	return s.recipeRepo.Delete(ctx, id)
}

func (s *RecipeService) UpdateWithIngredients(ctx context.Context, recipe *model.Recipe, ingredients []model.RecipeIngredient) error {
	tx, err := s.recipeRepo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Обновляем рецепт
	if err := s.recipeRepo.UpdateWithTx(ctx, tx, recipe); err != nil {
		return err
	}

	// Удаляем старые ингредиенты
	if err := s.recipeIngrRepo.DeleteByRecipeIDWithTx(ctx, tx, recipe.ID); err != nil {
		return err
	}

	// Добавляем новые ингредиенты
	for _, ing := range ingredients {
		ing.RecipeID = recipe.ID
		if err := s.recipeIngrRepo.AddWithTx(ctx, tx, &ing); err != nil {
			return err
		}
	}

	return tx.Commit()
}

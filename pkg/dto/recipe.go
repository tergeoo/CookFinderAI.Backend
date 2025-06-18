package dto

import (
	"CookFinder.Backend/internal/model"
	"time"
)

type RecipeResponse struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	PrepTimeMin int          `json:"prep_time_min"`
	CookTimeMin int          `json:"cook_time_min"`
	Method      string       `json:"method"`
	ImageURL    string       `json:"image_url"`
	CreatedAt   time.Time    `json:"created_at"`
	Category    *Category    `json:"category"`
	Ingredients []Ingredient `json:"ingredients"`
}

type RecipeRequest struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	CategoryID    string   `json:"category_id"`
	PrepTimeMin   int      `json:"prep_time_min"`
	CookTimeMin   int      `json:"cook_time_min"`
	Method        string   `json:"method"`
	ImageURL      string   `json:"image_url"`
	IngredientIDs []string `json:"ingredient_ids"`
}

func NewRecipeResponseFromModel(recipe *model.RecipeCategoryIngredients) *RecipeResponse {
	ingredients := make([]Ingredient, len(recipe.Ingredients))
	for i, ing := range recipe.Ingredients {
		ingredients[i] = *NewIngredientFromModel(&ing)
	}

	category := NewCategoryFromModel(&recipe.Category)

	return &RecipeResponse{
		ID:          recipe.Recipe.ID,
		Title:       recipe.Recipe.Title,
		PrepTimeMin: recipe.Recipe.PrepTimeMin,
		CookTimeMin: recipe.Recipe.CookTimeMin,
		Method:      recipe.Recipe.Method,
		CreatedAt:   recipe.Recipe.CreatedAt,
		ImageURL:    recipe.Recipe.ImageURL,
		Category:    category,
		Ingredients: ingredients,
	}
}

package dto

import (
	"CookFinder.Backend/internal/model"
	"time"
)

type RecipeResponse struct {
	ID          string                     `json:"id"`
	Title       string                     `json:"title"`
	PrepTimeMin int                        `json:"prep_time_min"`
	CookTimeMin int                        `json:"cook_time_min"`
	Method      string                     `json:"method"`
	ImageURL    string                     `json:"image_url"`
	Energy      int                        `json:"energy"`
	Fat         float64                    `json:"fat"`
	Protein     float64                    `json:"protein"`
	CreatedAt   time.Time                  `json:"created_at"`
	Category    *Category                  `json:"category"`
	Ingredients []RecipeIngredientResponse `json:"ingredients"`
}

type RecipeRequest struct {
	Title       string                    `json:"title"`
	CategoryID  string                    `json:"category_id"`
	PrepTimeMin int                       `json:"prep_time_min"`
	CookTimeMin int                       `json:"cook_time_min"`
	Energy      int                       `json:"energy"`
	Fat         float64                   `json:"fat"`
	Protein     float64                   `json:"protein"`
	Method      string                    `json:"method"`
	ImageURL    string                    `json:"image_url"`
	Ingredients []RecipeIngredientRequest `json:"ingredients"`
}

func NewRecipeResponseFromModel(recipe *model.RecipeCategoryIngredients) *RecipeResponse {
	ingredients := make([]RecipeIngredientResponse, len(recipe.Ingredients))
	for i, ing := range recipe.Ingredients {
		ingredients[i] = RecipeIngredientResponse{
			ID:     ing.ID,
			Name:   ing.Name,
			Amount: ing.Amount,
			Unit:   ing.Unit,
			Image:  ing.ImageURL,
		}
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
		Energy:      recipe.Recipe.Energy,
		Fat:         recipe.Recipe.Fat,
		Protein:     recipe.Recipe.Protein,
		Category:    category,
		Ingredients: ingredients,
	}
}

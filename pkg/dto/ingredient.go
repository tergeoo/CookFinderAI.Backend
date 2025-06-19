package dto

import "CookFinder.Backend/internal/model"

type IngredientRequest struct {
	Name     string `json:"name"`
	ImageUrl string `json:"image_url"`
}

type IngredientResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ImageUrl string `json:"image_url"`
}

func NewIngredientFromModel(ingredient *model.Ingredient) *IngredientResponse {
	return &IngredientResponse{
		ID:       ingredient.ID,
		Name:     ingredient.Name,
		ImageUrl: ingredient.ImageUrl,
	}
}

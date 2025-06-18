package dto

import "CookFinder.Backend/internal/model"

type Ingredient struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewIngredientFromModel(ingredient *model.Ingredient) *Ingredient {
	return &Ingredient{
		ID:   ingredient.ID,
		Name: ingredient.Name,
	}
}

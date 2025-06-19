package dto

type RecipeIngredientResponse struct {
	ID     string `json:"id"` // ingredient_id
	Name   string `json:"name"`
	Amount int    `json:"amount"`
	Unit   string `json:"unit"`
	Image  string `json:"image_url"`
}
type RecipeIngredientRequest struct {
	ID     string `json:"id"` // ingredient_id
	Amount int    `json:"amount"`
	Unit   string `json:"unit"`
}

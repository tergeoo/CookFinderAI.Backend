package model

type RecipeIngredient struct {
	RecipeID     string `db:"recipe_id" json:"recipe_id"`
	IngredientID string `db:"ingredient_id" json:"ingredient_id"`
	Amount       int    `db:"amount" json:"amount"`
	Unit         string `db:"unit" json:"unit"` // g, ml, pcs
}

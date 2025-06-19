package model

type RecipeCategoryIngredients struct {
	Recipe      Recipe
	Category    Category
	Ingredients []IngredientWithAmount
}

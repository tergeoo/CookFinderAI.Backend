package model

type IngredientWithAmount struct {
	ID       string `db:"id"`        // ingredients.id
	Name     string `db:"name"`      // ingredients.name
	ImageURL string `db:"image_url"` // ingredients.image_url
	Amount   int    `db:"amount"`    // recipe_ingredients.amount
	Unit     string `db:"unit"`      // recipe_ingredients.unit
}

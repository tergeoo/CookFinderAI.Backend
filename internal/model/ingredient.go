package model

type Ingredient struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	ImageUrl string `db:"image_url"`
}

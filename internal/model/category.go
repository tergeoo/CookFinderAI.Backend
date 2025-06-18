package model

type Category struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	ImageURL string `db:"image_url"`
}

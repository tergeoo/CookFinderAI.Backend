package dto

import "CookFinder.Backend/internal/model"

type Category struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

func NewCategoryFromModel(category *model.Category) *Category {
	return &Category{
		ID:       category.ID,
		Name:     category.Name,
		ImageURL: category.ImageUrl,
	}
}

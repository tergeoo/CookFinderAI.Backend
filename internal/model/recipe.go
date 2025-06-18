package model

import (
	"github.com/lib/pq"
	"time"
)

type Recipe struct {
	ID            string         `db:"id"`
	Title         string         `db:"title"`
	CategoryID    string         `db:"category_id"`
	PrepTimeMin   int            `db:"prep_time_min"`
	CookTimeMin   int            `db:"cook_time_min"`
	Method        string         `db:"method"`
	CreatedAt     time.Time      `db:"created_at"`
	ImageURL      string         `db:"image_url"`
	IngredientIDs pq.StringArray `db:"ingredient_ids"` // []string
}

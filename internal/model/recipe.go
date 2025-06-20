package model

import (
	"time"
)

type Recipe struct {
	ID          string    `db:"id"`
	Title       string    `db:"title"`
	CategoryID  string    `db:"category_id"`
	PrepTimeMin int       `db:"prep_time_min"`
	CookTimeMin int       `db:"cook_time_min"`
	Method      string    `db:"method"`
	Energy      int       `db:"energy"`
	Fat         float64   `db:"fat"`
	Protein     float64   `db:"protein"`
	CreatedAt   time.Time `db:"created_at"`
	ImageURL    string    `db:"image_url"`
}

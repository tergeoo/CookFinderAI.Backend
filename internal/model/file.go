package model

type File struct {
	ID   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Path string `db:"path" json:"path"`
}

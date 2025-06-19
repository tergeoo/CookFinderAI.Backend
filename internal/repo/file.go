package repo

import (
	"CookFinder.Backend/internal/model"
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type FileRepository struct {
	db *sqlx.DB
	sb squirrel.StatementBuilderType
}

func NewFileRepository(db *sqlx.DB) *FileRepository {
	return &FileRepository{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *FileRepository) Create(ctx context.Context, file *model.File) error {
	query, args, err := r.sb.Insert("files").
		Columns("id", "name", "path").
		Values(file.ID, file.Name, file.Path).
		Suffix("ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, path = EXCLUDED.path").
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *FileRepository) GetAll(ctx context.Context) ([]model.File, error) {
	query, args, err := r.sb.Select("*").
		From("files").
		OrderBy("name").
		ToSql()
	if err != nil {
		return nil, err
	}

	var files []model.File
	err = r.db.SelectContext(ctx, &files, query, args...)
	return files, err
}

func (r *FileRepository) GetByID(ctx context.Context, id string) (*model.File, error) {
	query, args, err := r.sb.Select("*").
		From("files").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var file model.File
	if err := r.db.GetContext(ctx, &file, query, args...); err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *FileRepository) Delete(ctx context.Context, id string) error {
	query, args, err := r.sb.
		Delete("files").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

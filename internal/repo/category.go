package repo

import (
	"CookFinder.Backend/internal/model"
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type CategoryRepository struct {
	db *sqlx.DB
	sb squirrel.StatementBuilderType
}

func NewCategoryRepository(db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (it *CategoryRepository) Create(ctx context.Context, category *model.Category) error {
	query, args, err := it.sb.Insert("recipe_categories").
		Columns("id", "name", "image_url").
		Values(category.ID, category.Name, category.ImageUrl).
		Suffix("ON CONFLICT (name) DO UPDATE SET image_url = EXCLUDED.image_url").
		ToSql()
	if err != nil {
		return err
	}

	_, err = it.db.ExecContext(ctx, query, args...)
	return err
}

func (it *CategoryRepository) GetAll(ctx context.Context) ([]model.Category, error) {
	query, args, err := it.sb.Select("*").
		From("recipe_categories").
		OrderBy("name").
		ToSql()
	if err != nil {
		return nil, err
	}

	var categories []model.Category
	err = it.db.SelectContext(ctx, &categories, query, args...)
	return categories, err
}

func (it *CategoryRepository) GetByID(ctx context.Context, id string) (*model.Category, error) {
	query, args, err := it.sb.Select("*").
		From("recipe_categories").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var category model.Category
	if err := it.db.GetContext(ctx, &category, query, args...); err != nil {
		return nil, err
	}
	return &category, nil
}

func (it *CategoryRepository) Delete(ctx context.Context, id string) error {
	query, args, err := it.sb.
		Delete("recipe_categories").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = it.db.ExecContext(ctx, query, args...)
	return err
}

func (it *CategoryRepository) Update(ctx context.Context, category *model.Category) error {
	query, args, err := it.sb.Update("ingredients").
		Set("name", category.Name).
		Set("image_url", category.ImageUrl).
		Where(squirrel.Eq{"id": category.ID}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = it.db.ExecContext(ctx, query, args...)
	return err
}

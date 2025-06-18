package repository

import (
	"CookFinder.Backend/internal/model"
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type IngredientRepository struct {
	db *sqlx.DB
	sb squirrel.StatementBuilderType
}

func NewIngredientRepository(db *sqlx.DB) *IngredientRepository {
	return &IngredientRepository{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (it *IngredientRepository) Create(ctx context.Context, ingredient *model.Ingredient) error {
	query, args, err := it.sb.Insert("ingredients").
		Columns("id", "name").
		Values(ingredient.ID, ingredient.Name).
		Suffix("ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name").
		ToSql()
	if err != nil {
		return err
	}
	_, err = it.db.ExecContext(ctx, query, args...)
	return err
}

func (it *IngredientRepository) GetByID(ctx context.Context, id string) (*model.Ingredient, error) {
	query, args, err := it.sb.Select("*").From("ingredients").Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, err
	}
	var ingredient model.Ingredient
	if err := it.db.GetContext(ctx, &ingredient, query, args...); err != nil {
		return nil, err
	}
	return &ingredient, nil
}

func (it *IngredientRepository) GetAllByID(ctx context.Context, id string) ([]model.Ingredient, error) {
	query, args, err := it.sb.Select("*").
		From("ingredients").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}
	var ingredients []model.Ingredient
	err = it.db.SelectContext(ctx, &ingredients, query, args...)
	return ingredients, err
}

func (it *IngredientRepository) GetAll(ctx context.Context) ([]model.Ingredient, error) {
	query, args, err := it.sb.Select("*").
		From("ingredients").
		OrderBy("name").
		ToSql()
	if err != nil {
		return nil, err
	}
	var ingredients []model.Ingredient
	err = it.db.SelectContext(ctx, &ingredients, query, args...)
	return ingredients, err
}

func (it *IngredientRepository) Delete(ctx context.Context, id string) error {
	query, args, err := it.sb.Delete("ingredients").Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return err
	}
	_, err = it.db.ExecContext(ctx, query, args...)
	return err
}

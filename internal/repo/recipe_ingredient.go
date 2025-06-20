package repo

import (
	"CookFinder.Backend/internal/model"
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type RecipeIngredientRepository struct {
	db *sqlx.DB
	sq squirrel.StatementBuilderType
}

func NewRecipeIngredientRepository(db *sqlx.DB) *RecipeIngredientRepository {
	return &RecipeIngredientRepository{
		db: db,
		sq: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (it *RecipeIngredientRepository) Add(ctx context.Context, ri *model.RecipeIngredient) error {
	query, args, err := it.sq.Insert("recipe_ingredients").
		Columns("recipe_id", "ingredient_id", "amount", "unit").
		Values(ri.RecipeID, ri.IngredientID, ri.Amount, ri.Unit).
		Suffix("ON CONFLICT (recipe_id, ingredient_id) DO UPDATE SET amount = EXCLUDED.amount, unit = EXCLUDED.unit").
		ToSql()
	if err != nil {
		return err
	}
	_, err = it.db.ExecContext(ctx, query, args...)
	return err
}

// ✅ Новый метод AddWithTx
func (it *RecipeIngredientRepository) AddWithTx(ctx context.Context, tx *sqlx.Tx, ri *model.RecipeIngredient) error {
	query, args, err := it.sq.Insert("recipe_ingredients").
		Columns("recipe_id", "ingredient_id", "amount", "unit").
		Values(ri.RecipeID, ri.IngredientID, ri.Amount, ri.Unit).
		Suffix("ON CONFLICT (recipe_id, ingredient_id) DO UPDATE SET amount = EXCLUDED.amount, unit = EXCLUDED.unit").
		ToSql()
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, query, args...)
	return err
}

func (it *RecipeIngredientRepository) GetByRecipeID(ctx context.Context, recipeID string) ([]model.RecipeIngredient, error) {
	query, args, err := it.sq.Select("*").
		From("recipe_ingredients").
		Where(squirrel.Eq{"recipe_id": recipeID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var result []model.RecipeIngredient
	err = it.db.SelectContext(ctx, &result, query, args...)
	return result, err
}

func (it *RecipeIngredientRepository) DeleteByRecipeID(ctx context.Context, recipeID string) error {
	query, args, err := it.sq.Delete("recipe_ingredients").
		Where(squirrel.Eq{"recipe_id": recipeID}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = it.db.ExecContext(ctx, query, args...)
	return err
}

func (it *RecipeIngredientRepository) DeleteByRecipeIDWithTx(ctx context.Context, tx *sqlx.Tx, recipeID string) error {
	queryBuilder := it.sq.
		Delete("recipe_ingredients").
		Where(squirrel.Eq{"recipe_id": recipeID}).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	return err
}

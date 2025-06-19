package repo

import (
	"CookFinder.Backend/internal/model"
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type RecipeIngredientRepository struct {
	db *sqlx.DB
	sb squirrel.StatementBuilderType
}

func NewRecipeIngredientRepository(db *sqlx.DB) *RecipeIngredientRepository {
	return &RecipeIngredientRepository{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *RecipeIngredientRepository) Add(ctx context.Context, ri *model.RecipeIngredient) error {
	query, args, err := r.sb.Insert("recipe_ingredients").
		Columns("recipe_id", "ingredient_id", "amount", "unit").
		Values(ri.RecipeID, ri.IngredientID, ri.Amount, ri.Unit).
		Suffix("ON CONFLICT (recipe_id, ingredient_id) DO UPDATE SET amount = EXCLUDED.amount, unit = EXCLUDED.unit").
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

// ✅ Новый метод AddWithTx
func (r *RecipeIngredientRepository) AddWithTx(ctx context.Context, tx *sqlx.Tx, ri *model.RecipeIngredient) error {
	query, args, err := r.sb.Insert("recipe_ingredients").
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

func (r *RecipeIngredientRepository) GetByRecipeID(ctx context.Context, recipeID string) ([]model.RecipeIngredient, error) {
	query, args, err := r.sb.Select("*").
		From("recipe_ingredients").
		Where(squirrel.Eq{"recipe_id": recipeID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var result []model.RecipeIngredient
	err = r.db.SelectContext(ctx, &result, query, args...)
	return result, err
}

func (r *RecipeIngredientRepository) DeleteByRecipeID(ctx context.Context, recipeID string) error {
	query, args, err := r.sb.Delete("recipe_ingredients").
		Where(squirrel.Eq{"recipe_id": recipeID}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

package repo

import (
	"CookFinder.Backend/internal/model"
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RecipeRepository struct {
	db *sqlx.DB
	sb squirrel.StatementBuilderType
}

func NewRecipeRepository(db *sqlx.DB) *RecipeRepository {
	return &RecipeRepository{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *RecipeRepository) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return r.db.BeginTxx(ctx, nil)
}

func (r *RecipeRepository) Create(ctx context.Context, recipe *model.Recipe) error {
	query, args, err := r.sb.
		Insert("recipes").
		Columns("id", "title", "category_id", "prep_time_min", "cook_time_min", "method", "created_at", "image_url").
		Values(recipe.ID, recipe.Title, recipe.CategoryID, recipe.PrepTimeMin, recipe.CookTimeMin, recipe.Method, recipe.CreatedAt, recipe.ImageURL).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *RecipeRepository) CreateWithTx(ctx context.Context, tx *sqlx.Tx, recipe *model.Recipe) error {
	query, args, err := r.sb.
		Insert("recipes").
		Columns("id", "title", "category_id", "prep_time_min", "cook_time_min", "method", "created_at", "image_url").
		Values(recipe.ID, recipe.Title, recipe.CategoryID, recipe.PrepTimeMin, recipe.CookTimeMin, recipe.Method, recipe.CreatedAt, recipe.ImageURL).
		ToSql()
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, query, args...)
	return err
}

func (r *RecipeRepository) GetByID(ctx context.Context, id string) (*model.RecipeCategoryIngredients, error) {
	query, args, err := r.sb.
		Select(
			"r.id", "r.title", "r.category_id", "r.prep_time_min", "r.cook_time_min", "r.method", "r.created_at", "r.image_url",
			"c.id AS category_id", "c.name AS category_name", "c.image_url AS category_image_url",
		).
		From("recipes r").
		Join("recipe_categories c ON r.category_id = c.id").
		Where(squirrel.Eq{"r.id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var row struct {
		model.Recipe
		CategoryID       string `db:"category_id"`
		CategoryName     string `db:"category_name"`
		CategoryImageURL string `db:"category_image_url"`
	}
	if err := r.db.GetContext(ctx, &row, query, args...); err != nil {
		return nil, err
	}

	ingredients, err := r.getIngredientsByRecipeID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &model.RecipeCategoryIngredients{
		Recipe: row.Recipe,
		Category: model.Category{
			ID:       row.CategoryID,
			Name:     row.CategoryName,
			ImageUrl: row.CategoryImageURL,
		},
		Ingredients: ingredients,
	}, nil
}

func (r *RecipeRepository) GetAll(ctx context.Context) ([]model.RecipeCategoryIngredients, error) {
	query, args, err := r.sb.
		Select(
			"r.id", "r.title", "r.category_id", "r.prep_time_min", "r.cook_time_min", "r.method", "r.created_at", "r.image_url",
			"c.id AS category_id", "c.name AS category_name", "c.image_url AS category_image_url",
		).
		From("recipes r").
		Join("recipe_categories c ON r.category_id = c.id").
		OrderBy("r.created_at DESC").
		ToSql()
	if err != nil {
		return nil, err
	}

	var rows []struct {
		model.Recipe
		CategoryID       string `db:"category_id"`
		CategoryName     string `db:"category_name"`
		CategoryImageURL string `db:"category_image_url"`
	}
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, err
	}

	var result []model.RecipeCategoryIngredients
	for _, row := range rows {
		ingredients, err := r.getIngredientsByRecipeID(ctx, row.ID)
		if err != nil {
			return nil, err
		}

		result = append(result, model.RecipeCategoryIngredients{
			Recipe: row.Recipe,
			Category: model.Category{
				ID:       row.CategoryID,
				Name:     row.CategoryName,
				ImageUrl: row.CategoryImageURL,
			},
			Ingredients: ingredients,
		})
	}

	return result, nil
}

func (r *RecipeRepository) getIngredientsByRecipeID(ctx context.Context, recipeID string) ([]model.IngredientWithAmount, error) {
	query, args, err := r.sb.
		Select("i.id", "i.name", "i.image_url", "ri.amount", "ri.unit").
		From("recipe_ingredients ri").
		Join("ingredients i ON i.id = ri.ingredient_id").
		Where(squirrel.Eq{"ri.recipe_id": recipeID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var ingredients []model.IngredientWithAmount
	if err := r.db.SelectContext(ctx, &ingredients, query, args...); err != nil {
		return nil, err
	}
	return ingredients, nil
}

func (r *RecipeRepository) Update(ctx context.Context, recipe *model.Recipe) error {
	query, args, err := r.sb.Update("recipes").
		Set("title", recipe.Title).
		Set("category_id", recipe.CategoryID).
		Set("prep_time_min", recipe.PrepTimeMin).
		Set("cook_time_min", recipe.CookTimeMin).
		Set("method", recipe.Method).
		Set("image_url", recipe.ImageURL).
		Where(squirrel.Eq{"id": recipe.ID}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *RecipeRepository) Delete(ctx context.Context, id string) error {
	query, args, err := r.sb.Delete("recipes").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *RecipeRepository) BatchInsert(ctx context.Context, recipes []model.Recipe) error {
	q := r.sb.Insert("recipes").
		Columns("id", "title", "category_id", "prep_time_min", "cook_time_min", "method", "created_at", "image_url")

	for _, rec := range recipes {
		if rec.ID == "" {
			rec.ID = uuid.New().String()
		}
		q = q.Values(rec.ID, rec.Title, rec.CategoryID, rec.PrepTimeMin, rec.CookTimeMin, rec.Method, time.Now(), rec.ImageURL)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *RecipeRepository) BatchDelete(ctx context.Context, ids []string) error {
	query, args, err := r.sb.Delete("recipes").
		Where(squirrel.Eq{"id": ids}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

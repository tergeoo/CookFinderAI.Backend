package repository

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

func (it *RecipeRepository) Create(ctx context.Context, recipe *model.Recipe) error {
	query, args, err := it.sb.
		Insert("recipes").
		Columns("id", "title", "category_id", "prep_time_min", "cook_time_min", "method", "created_at", "ingredient_ids", "image_url").
		Values(recipe.ID, recipe.Title, recipe.CategoryID, recipe.PrepTimeMin, recipe.CookTimeMin, recipe.Method, time.Now(), recipe.IngredientIDs, recipe.ImageURL).
		ToSql()
	if err != nil {
		return err
	}

	_, err = it.db.ExecContext(ctx, query, args...)
	return err
}

func (it *RecipeRepository) GetByID(ctx context.Context, id string) (*model.RecipeCategoryIngredients, error) {
	query, args, err := it.sb.
		Select(
			"r.id",
			"r.title",
			"r.category_id",
			"r.prep_time_min",
			"r.cook_time_min",
			"r.method",
			"r.created_at",
			"r.image_url",
			"r.ingredient_ids",

			"c.id AS category_id",
			"c.name AS category_name",
			"c.image_url AS category_image_url",
		).
		From("recipes r").
		Join("categories c ON r.category_id = c.id").
		Where(squirrel.Eq{"r.id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var joined struct {
		model.Recipe
		CategoryID       string `db:"category_id"`
		CategoryName     string `db:"category_name"`
		CategoryImageURL string `db:"category_image_url"`
	}

	if err := it.db.GetContext(ctx, &joined, query, args...); err != nil {
		return nil, err
	}

	// Загружаем ингредиенты по ID
	ingQuery, ingArgs, _ := it.sb.
		Select("id", "name").
		From("ingredients").
		Where(squirrel.Eq{"id": joined.IngredientIDs}).
		ToSql()

	var ingredients []model.Ingredient
	if err := it.db.SelectContext(ctx, &ingredients, ingQuery, ingArgs...); err != nil {
		return nil, err
	}

	return &model.RecipeCategoryIngredients{
		Recipe: joined.Recipe,
		Category: model.Category{
			ID:       joined.CategoryID,
			Name:     joined.CategoryName,
			ImageURL: joined.CategoryImageURL,
		},
		Ingredients: ingredients,
	}, nil
}

func (it *RecipeRepository) GetAll(ctx context.Context) ([]model.RecipeCategoryIngredients, error) {
	query, args, err := it.sb.
		Select(
			"r.id", "r.title", "r.category_id", "r.prep_time_min", "r.cook_time_min", "r.method",
			"r.created_at", "r.image_url", "r.ingredient_ids",
			"c.id AS category_id", "c.name AS category_name", "c.image_url AS category_image_url",
		).
		From("recipes r").
		Join("categories c ON r.category_id = c.id").
		OrderBy("r.created_at DESC").
		ToSql()
	if err != nil {
		return nil, err
	}

	type joinedRow struct {
		model.Recipe
		CategoryID       string `db:"category_id"`
		CategoryName     string `db:"category_name"`
		CategoryImageURL string `db:"category_image_url"`
	}

	var rows []joinedRow
	if err := it.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, err
	}

	var result []model.RecipeCategoryIngredients
	for _, row := range rows {
		// загрузка ингредиентов
		ingQuery, ingArgs, _ := it.sb.
			Select("id", "name").
			From("ingredients").
			Where(squirrel.Eq{"id": row.IngredientIDs}).
			ToSql()

		var ingredients []model.Ingredient
		if err := it.db.SelectContext(ctx, &ingredients, ingQuery, ingArgs...); err != nil {
			return nil, err
		}

		result = append(result, model.RecipeCategoryIngredients{
			Recipe: row.Recipe,
			Category: model.Category{
				ID:       row.CategoryID,
				Name:     row.CategoryName,
				ImageURL: row.CategoryImageURL,
			},
			Ingredients: ingredients,
		})
	}

	return result, nil
}

func (it *RecipeRepository) Update(ctx context.Context, recipe *model.Recipe) error {
	query, args, err := it.sb.Update("recipes").
		Set("title", recipe.Title).
		Set("category_id", recipe.CategoryID).
		Set("prep_time_min", recipe.PrepTimeMin).
		Set("cook_time_min", recipe.CookTimeMin).
		Set("method", recipe.Method).
		Set("ingredient_ids", recipe.IngredientIDs).
		Set("image_url", recipe.ImageURL).
		Where(squirrel.Eq{"id": recipe.ID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = it.db.ExecContext(ctx, query, args...)
	return err
}

func (it *RecipeRepository) Delete(ctx context.Context, id string) error {
	query, args, err := it.sb.Delete("recipes").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = it.db.ExecContext(ctx, query, args...)
	return err
}

func (it *RecipeRepository) BatchInsert(ctx context.Context, recipes []model.Recipe) error {
	q := it.sb.Insert("recipes").
		Columns("id", "title", "category_id", "prep_time_min", "cook_time_min", "method", "created_at", "ingredient_ids", "image_url")

	for _, rec := range recipes {
		if rec.ID == "" {
			rec.ID = uuid.New().String()
		}
		q = q.Values(rec.ID, rec.Title, rec.CategoryID, rec.PrepTimeMin, rec.CookTimeMin, rec.Method, time.Now(), rec.IngredientIDs, rec.ImageURL)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return err
	}

	_, err = it.db.ExecContext(ctx, query, args...)
	return err
}

func (it *RecipeRepository) BatchDelete(ctx context.Context, ids []string) error {
	query, args, err := it.sb.Delete("recipes").
		Where(squirrel.Eq{"id": ids}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = it.db.ExecContext(ctx, query, args...)
	return err
}

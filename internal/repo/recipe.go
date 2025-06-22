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
	sq squirrel.StatementBuilderType
}

func NewRecipeRepository(db *sqlx.DB) *RecipeRepository {
	return &RecipeRepository{
		db: db,
		sq: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (it *RecipeRepository) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return it.db.BeginTxx(ctx, nil)
}

func (it *RecipeRepository) Create(ctx context.Context, recipe *model.Recipe) error {
	query, args, err := it.sq.
		Insert("recipes").
		Columns("id", "title", "category_id", "prep_time_min", "cook_time_min", "method", "created_at", "image_url", "energy", "fat", "protein").
		Values(recipe.ID, recipe.Title, recipe.CategoryID, recipe.PrepTimeMin, recipe.CookTimeMin, recipe.Method, recipe.CreatedAt, recipe.ImageURL, recipe.Energy, recipe.Fat, recipe.Protein).
		ToSql()
	if err != nil {
		return err
	}
	_, err = it.db.ExecContext(ctx, query, args...)
	return err
}

func (it *RecipeRepository) CreateWithTx(ctx context.Context, tx *sqlx.Tx, recipe *model.Recipe) error {
	query, args, err := it.sq.
		Insert("recipes").
		Columns("id", "title", "category_id", "prep_time_min", "cook_time_min", "method", "created_at", "image_url", "energy", "fat", "protein").
		Values(recipe.ID, recipe.Title, recipe.CategoryID, recipe.PrepTimeMin, recipe.CookTimeMin, recipe.Method, recipe.CreatedAt, recipe.ImageURL, recipe.Energy, recipe.Fat, recipe.Protein).
		ToSql()
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, query, args...)
	return err
}

func (it *RecipeRepository) GetByID(ctx context.Context, id string) (*model.RecipeCategoryIngredients, error) {
	query, args, err := it.sq.
		Select(
			"it.id", "it.title", "it.category_id", "it.prep_time_min", "it.cook_time_min", "it.method", "it.created_at", "it.image_url", "it.energy", "it.fat", "it.protein",
			"c.id AS category_id", "c.name AS category_name", "c.image_url AS category_image_url",
		).
		From("recipes it").
		Join("recipe_categories c ON it.category_id = c.id").
		Where(squirrel.Eq{"it.id": id}).
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
	if err := it.db.GetContext(ctx, &row, query, args...); err != nil {
		return nil, err
	}

	ingredients, err := it.getIngredientsByRecipeID(ctx, id)
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

func (it *RecipeRepository) GetAll(ctx context.Context, search, categoryID string) ([]model.RecipeCategoryIngredients, error) {
	builder := it.sq.
		Select(
			"DISTINCT it.id", "it.title", "it.category_id", "it.prep_time_min", "it.cook_time_min", "it.method", "it.created_at", "it.image_url", "it.energy", "it.fat", "it.protein",
			"c.id AS category_id", "c.name AS category_name", "c.image_url AS category_image_url",
		).
		From("recipes it").
		Join("recipe_categories c ON it.category_id = c.id").
		LeftJoin("recipe_ingredients ri ON ri.recipe_id = it.id").
		LeftJoin("ingredients i ON i.id = ri.ingredient_id").
		OrderBy("it.created_at DESC")

	// Фильтрация по названию рецепта и ингредиентам
	if search != "" {
		builder = builder.Where(
			squirrel.Or{
				squirrel.Expr("LOWER(it.title) LIKE LOWER(?)", "%"+search+"%"),
				squirrel.Expr("LOWER(i.name) LIKE LOWER(?)", "%"+search+"%"),
			},
		)
	}

	// Фильтрация по категории
	if categoryID != "" {
		builder = builder.Where("it.category_id = ?", categoryID)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var rows []struct {
		model.Recipe
		CategoryID       string `db:"category_id"`
		CategoryName     string `db:"category_name"`
		CategoryImageURL string `db:"category_image_url"`
	}
	if err := it.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, err
	}

	var result []model.RecipeCategoryIngredients
	for _, row := range rows {
		ingredients, err := it.getIngredientsByRecipeID(ctx, row.ID)
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

func (it *RecipeRepository) getIngredientsByRecipeID(ctx context.Context, recipeID string) ([]model.IngredientWithAmount, error) {
	query, args, err := it.sq.
		Select("i.id", "i.name", "i.image_url", "ri.amount", "ri.unit").
		From("recipe_ingredients ri").
		Join("ingredients i ON i.id = ri.ingredient_id").
		Where(squirrel.Eq{"ri.recipe_id": recipeID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var ingredients []model.IngredientWithAmount
	if err := it.db.SelectContext(ctx, &ingredients, query, args...); err != nil {
		return nil, err
	}
	return ingredients, nil
}

func (it *RecipeRepository) Update(ctx context.Context, recipe *model.Recipe) error {
	query, args, err := it.sq.Update("recipes").
		Set("title", recipe.Title).
		Set("category_id", recipe.CategoryID).
		Set("prep_time_min", recipe.PrepTimeMin).
		Set("cook_time_min", recipe.CookTimeMin).
		Set("method", recipe.Method).
		Set("image_url", recipe.ImageURL).
		Set("energy", recipe.Energy).
		Set("fat", recipe.Fat).
		Set("protein", recipe.Protein).
		Where(squirrel.Eq{"id": recipe.ID}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = it.db.ExecContext(ctx, query, args...)
	return err
}

func (it *RecipeRepository) Delete(ctx context.Context, id string) error {
	query, args, err := it.sq.Delete("recipes").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = it.db.ExecContext(ctx, query, args...)
	return err
}

func (it *RecipeRepository) BatchInsert(ctx context.Context, recipes []model.Recipe) error {
	q := it.sq.Insert("recipes").
		Columns("id", "title", "category_id", "prep_time_min", "cook_time_min", "method", "created_at", "image_url", "energy", "fat", "protein")

	for _, rec := range recipes {
		if rec.ID == "" {
			rec.ID = uuid.New().String()
		}
		q = q.Values(rec.ID, rec.Title, rec.CategoryID, rec.PrepTimeMin, rec.CookTimeMin, rec.Method, time.Now(), rec.ImageURL, rec.Energy, rec.Fat, rec.Protein)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return err
	}
	_, err = it.db.ExecContext(ctx, query, args...)
	return err
}

func (it *RecipeRepository) BatchDelete(ctx context.Context, ids []string) error {
	query, args, err := it.sq.Delete("recipes").
		Where(squirrel.Eq{"id": ids}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = it.db.ExecContext(ctx, query, args...)
	return err
}

func (it *RecipeRepository) UpdateWithTx(ctx context.Context, tx *sqlx.Tx, recipe *model.Recipe) error {
	queryBuilder := it.sq.
		Update("recipes").
		Set("title", recipe.Title).
		Set("category_id", recipe.CategoryID).
		Set("prep_time_min", recipe.PrepTimeMin).
		Set("cook_time_min", recipe.CookTimeMin).
		Set("method", recipe.Method).
		Set("energy", recipe.Energy).
		Set("fat", recipe.Fat).
		Set("protein", recipe.Protein).
		Set("image_url", recipe.ImageURL).
		Where(squirrel.Eq{"id": recipe.ID})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	return err
}

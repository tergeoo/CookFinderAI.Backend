package handler

import (
	"CookFinder.Backend/internal/model"
	"CookFinder.Backend/internal/service"
	"CookFinder.Backend/pkg/dto"
	"CookFinder.Backend/pkg/uuid"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RecipeHandler struct {
	service *service.RecipeService
}

func NewRecipeHandler(r *gin.Engine, svc *service.RecipeService) {
	h := &RecipeHandler{service: svc}
	routes := r.Group("/recipes")
	{
		routes.POST("", h.Create)
		routes.GET("", h.GetAll)
		routes.GET(":id", h.GetByID)
		routes.PUT(":id", h.Update)
		routes.DELETE(":id", h.Delete)
	}
}

// GetAll godoc
// @Summary Get all recipes
// @Tags Recipes
// @Produce json
// @Success 200 {array} dto.RecipeResponse
// @Failure 500 {object} map[string]string
// @Router /recipes [get]
func (h *RecipeHandler) GetAll(c *gin.Context) {
	recipes, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	results := make([]dto.RecipeResponse, 0, len(recipes))
	for _, r := range recipes {
		results = append(results, *dto.NewRecipeResponseFromModel(&r))
	}

	c.JSON(http.StatusOK, results)
}

// GetByID godoc
// @Summary Get recipe by ID
// @Tags Recipes
// @Produce json
// @Param id path string true "Recipe ID"
// @Success 200 {object} dto.RecipeResponse
// @Failure 404 {object} map[string]string
// @Router /recipes/{id} [get]
func (h *RecipeHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	recipe, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	result := dto.NewRecipeResponseFromModel(recipe)
	c.JSON(http.StatusOK, result)
}

// Create godoc
// @Summary Create a new recipe with ingredients
// @Tags Recipes
// @Accept json
// @Produce json
// @Param recipe body dto.RecipeRequest true "Recipe data"
// @Success 201 {object} dto.RecipeResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /recipes [post]
func (h *RecipeHandler) Create(c *gin.Context) {
	var input dto.RecipeRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipe := &model.Recipe{
		ID:          uuid.V7().String(),
		Title:       input.Title,
		CategoryID:  input.CategoryID,
		PrepTimeMin: input.PrepTimeMin,
		CookTimeMin: input.CookTimeMin,
		Method:      input.Method,
		ImageURL:    input.ImageURL,
		Protein:     input.Protein,
		Fat:         input.Fat,
		Energy:      input.Energy,
		CreatedAt:   time.Now(),
	}

	ingredients := make([]model.RecipeIngredient, len(input.Ingredients))
	for i, ing := range input.Ingredients {
		ingredients[i] = model.RecipeIngredient{
			IngredientID: ing.ID,
			Amount:       ing.Amount,
			Unit:         ing.Unit,
		}
	}

	if err := h.service.CreateWithIngredients(c.Request.Context(), recipe, ingredients); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	created, err := h.service.GetByID(c.Request.Context(), recipe.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch created recipe"})
		return
	}

	c.JSON(http.StatusCreated, dto.NewRecipeResponseFromModel(created))
}

// Update godoc
// @Summary Update recipe by ID
// @Tags Recipes
// @Accept json
// @Produce json
// @Param id path string true "Recipe ID"
// @Param recipe body dto.RecipeRequest true "Recipe data"
// @Success 200 {object} dto.RecipeResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /recipes/{id} [put]
func (h *RecipeHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var input dto.RecipeRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверим, существует ли рецепт
	existing, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
		return
	}

	// Обновлённые данные рецепта
	updated := &model.Recipe{
		ID:          id,
		Title:       input.Title,
		CategoryID:  input.CategoryID,
		PrepTimeMin: input.PrepTimeMin,
		CookTimeMin: input.CookTimeMin,
		Method:      input.Method,
		ImageURL:    input.ImageURL,
		Protein:     input.Protein,
		Fat:         input.Fat,
		Energy:      input.Energy,
		CreatedAt:   existing.Recipe.CreatedAt,
	}

	// Новые ингредиенты
	ingredients := make([]model.RecipeIngredient, len(input.Ingredients))
	for i, ing := range input.Ingredients {
		ingredients[i] = model.RecipeIngredient{
			IngredientID: ing.ID,
			Amount:       ing.Amount,
			Unit:         ing.Unit,
		}
	}

	// Обновление рецепта и его ингредиентов
	if err := h.service.UpdateWithIngredients(c.Request.Context(), updated, ingredients); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Delete godoc
// @Summary Delete recipe by ID
// @Tags Recipes
// @Produce json
// @Param id path string true "Recipe ID"
// @Success 204 {string} string "No Content"
// @Failure 500 {object} map[string]string
// @Router /recipes/{id} [delete]
func (h *RecipeHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

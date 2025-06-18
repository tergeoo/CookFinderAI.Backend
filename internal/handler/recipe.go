package handler

import (
	"CookFinder.Backend/internal/model"
	"CookFinder.Backend/internal/service"
	"CookFinder.Backend/pkg/dto"
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
		routes.GET("", h.GetAll)
		routes.GET(":id", h.GetByID)
		routes.POST("", h.Create)
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
// @Param id path string true "RecipeResponse ID"
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
// @Summary Create a new recipe
// @Tags Recipes
// @Accept json
// @Produce json
// @Param recipe body dto.RecipeResponse true "RecipeResponse body"
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

	result := &model.Recipe{
		Title:         input.Title,
		CategoryID:    input.CategoryID,
		PrepTimeMin:   input.PrepTimeMin,
		CookTimeMin:   input.CookTimeMin,
		Method:        input.Method,
		CreatedAt:     time.Now(),
		ImageURL:      input.ImageURL,
		IngredientIDs: input.IngredientIDs,
	}

	if err := h.service.Create(c.Request.Context(), result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

// Update godoc
// @Summary Update recipe by ID
// @Tags Recipes
// @Accept json
// @Produce json
// @Param id path string true "RecipeResponse ID"
// @Param recipe body dto.RecipeResponse true "RecipeResponse body"
// @Success 200 {object} dto.RecipeResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /recipes/{id} [put]
func (h *RecipeHandler) Update(c *gin.Context) {
	var input dto.RecipeRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := &model.Recipe{
		ID:            c.Param("id"),
		Title:         input.Title,
		PrepTimeMin:   input.PrepTimeMin,
		CookTimeMin:   input.CookTimeMin,
		Method:        input.Method,
		ImageURL:      input.ImageURL,
		IngredientIDs: input.IngredientIDs,
	}

	if err := h.service.Update(c.Request.Context(), result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, input)
}

// Delete godoc
// @Summary Delete recipe by ID
// @Tags Recipes
// @Produce json
// @Param id path string true "RecipeResponse ID"
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

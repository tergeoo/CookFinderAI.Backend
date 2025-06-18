package handler

import (
	"CookFinder.Backend/internal/model"
	"CookFinder.Backend/internal/service"
	"CookFinder.Backend/pkg/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IngredientHandler struct {
	service *service.IngredientService
}

func NewIngredientHandler(r *gin.Engine, svc *service.IngredientService) {
	h := &IngredientHandler{service: svc}
	routes := r.Group("/ingredients")
	{
		routes.GET("", h.GetAll)
		routes.GET(":id", h.GetByID)
		routes.POST("", h.Create)
		routes.DELETE(":id", h.Delete)
	}
}

// GetAll godoc
// @Summary Get all ingredients
// @Tags IngredientIDs
// @Produce json
// @Success 200 {array} model.Ingredient
// @Failure 500 {object} map[string]string
// @Router /ingredients [get]
func (h *IngredientHandler) GetAll(c *gin.Context) {
	ingredients, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	results := make([]dto.Ingredient, 0, len(ingredients))

	for _, ing := range ingredients {
		results = append(results, *dto.NewIngredientFromModel(&ing))
	}

	c.JSON(http.StatusOK, results)
}

// GetByID godoc
// @Summary Get ingredient by ID
// @Tags IngredientIDs
// @Produce json
// @Param id path string true "Ingredient ID"
// @Success 200 {object} model.Ingredient
// @Failure 404 {object} map[string]string
// @Router /ingredients/{id} [get]
func (h *IngredientHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	ingredient, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	result := dto.NewIngredientFromModel(ingredient)

	c.JSON(http.StatusOK, result)
}

// Create godoc
// @Summary Create a new ingredient
// @Tags IngredientIDs
// @Accept json
// @Produce json
// @Param ingredient body model.Ingredient true "Ingredient body"
// @Success 201 {object} model.Ingredient
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /ingredients [post]
func (h *IngredientHandler) Create(c *gin.Context) {
	var input model.Ingredient
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := &model.Ingredient{
		Name: input.Name,
	}

	if err := h.service.Create(c.Request.Context(), result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

// Delete godoc
// @Summary Delete ingredient by ID
// @Tags IngredientIDs
// @Param id path string true "Ingredient ID"
// @Success 204
// @Failure 500 {object} map[string]string
// @Router /ingredients/{id} [delete]
func (h *IngredientHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

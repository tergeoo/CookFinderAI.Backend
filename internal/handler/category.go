package handler

import (
	"CookFinder.Backend/internal/model"
	"CookFinder.Backend/internal/service"
	"CookFinder.Backend/pkg/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(r *gin.Engine, svc *service.CategoryService) {
	h := &CategoryHandler{service: svc}
	routes := r.Group("/categories")
	{
		routes.GET("", h.GetAll)
		routes.GET(":id", h.GetByID)
		routes.POST("", h.Create)
		routes.DELETE(":id", h.Delete)
	}
}

// GetAll godoc
// @Summary GetAll all categories
// @Tags Categories
// @Produce json
// @Success 200 {array} dto.Category
// @Failure 500 {object} map[string]string
// @Router /categories [get]
func (h *CategoryHandler) GetAll(c *gin.Context) {
	categories, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	results := make([]dto.Category, 0, len(categories))

	for _, c := range categories {
		results = append(results, *dto.NewCategoryFromModel(&c))
	}

	c.JSON(http.StatusOK, results)
}

// GetByID godoc
// @Summary GetAll category by ID
// @Tags Categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} dto.Category
// @Failure 404 {object} map[string]string
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	category, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	res := dto.NewCategoryFromModel(category)
	c.JSON(http.StatusOK, res)
}

// Create godoc
// @Summary Create a new category
// @Tags Categories
// @Accept json
// @Produce json
// @Param category body dto.Category true "Category body"
// @Success 201 {object} dto.Category
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	var input dto.Category
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rc := &model.Category{
		Name:     input.Name,
		ImageURL: input.ImageURL,
	}

	if err := h.service.Create(c.Request.Context(), rc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, input)
}

// Delete godoc
// @Summary Delete category by ID
// @Tags Categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 204 {string} string "No Content"
// @Failure 500 {object} map[string]string
// @Router /categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

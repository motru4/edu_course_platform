package handler

import (
	"course2/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService *services.CategoryService
}

func NewCategoryHandler(router *gin.Engine, categoryService *services.CategoryService) {
	handler := &CategoryHandler{
		categoryService: categoryService,
	}

	categories := router.Group("/api/v1/edu/categories")
	{
		categories.GET("", handler.ListCategories)
	}
}

// @Summary Список категорий
// @Description Получить список всех категорий курсов
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {array} models.Category
// @Router /categories [get]
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	categories, err := h.categoryService.ListCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// controllers/category_controller.go
package controllers

import (
	"net/http"
	"news-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryController struct {
	DB *gorm.DB
}

func (c *CategoryController) CreateCategory(ctx *gin.Context) {
	var category models.Category
	if err := ctx.ShouldBindJSON(&category); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := c.DB.Create(&category)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, category)
}

func (c *CategoryController) GetCategories(ctx *gin.Context) {
	var categories []models.Category

	result := c.DB.Find(&categories)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, categories)
}

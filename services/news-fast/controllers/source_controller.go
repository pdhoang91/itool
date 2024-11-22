// controllers/source_controller.go
package controllers

import (
	"net/http"
	"news-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SourceController struct {
	DB *gorm.DB
}

func (c *SourceController) CreateSource(ctx *gin.Context) {
	var source models.Source
	if err := ctx.ShouldBindJSON(&source); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := c.DB.Create(&source)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, source)
}

func (c *SourceController) GetSources(ctx *gin.Context) {
	var sources []models.Source

	result := c.DB.Find(&sources)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, sources)
}

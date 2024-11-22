// controllers/news_controller.go
package controllers

import (
	"net/http"
	"strconv"

	"news-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NewsController struct {
	DB *gorm.DB
}

type PaginationResponse struct {
	Items       interface{} `json:"items"`
	CurrentPage int         `json:"current_page"`
	TotalPages  int         `json:"total_pages"`
	TotalItems  int64       `json:"total_items"`
	HasNext     bool        `json:"has_next"`
}

func (c *NewsController) GetNews(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("size", "20"))

	var news []models.News
	var totalItems int64

	// Get total count
	c.DB.Model(&models.News{}).Count(&totalItems)

	// Calculate total pages
	totalPages := int((totalItems + int64(pageSize) - 1) / int64(pageSize))

	// Get paginated news
	offset := (page - 1) * pageSize
	result := c.DB.
		Preload("Source").
		Preload("Category").
		Order("published_at desc").
		Offset(offset).
		Limit(pageSize).
		Find(&news)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	response := PaginationResponse{
		Items:       news,
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalItems:  totalItems,
		HasNext:     page < totalPages,
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *NewsController) GetNewsBySource(ctx *gin.Context) {
	sourceID := ctx.Param("id")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("size", "20"))

	var news []models.News
	var totalItems int64

	// Get total count for source
	c.DB.Model(&models.News{}).Where("source_id = ?", sourceID).Count(&totalItems)

	// Calculate total pages
	totalPages := int((totalItems + int64(pageSize) - 1) / int64(pageSize))

	// Get paginated news for source
	offset := (page - 1) * pageSize
	result := c.DB.
		Preload("Source").
		Preload("Category").
		Where("source_id = ?", sourceID).
		Order("published_at desc").
		Offset(offset).
		Limit(pageSize).
		Find(&news)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	response := PaginationResponse{
		Items:       news,
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalItems:  totalItems,
		HasNext:     page < totalPages,
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *NewsController) GetNewsDetail(ctx *gin.Context) {
	id := ctx.Param("id")

	var news models.News
	result := c.DB.
		Preload("Source").
		Preload("Category").
		First(&news, id)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
		return
	}

	ctx.JSON(http.StatusOK, news)
}

func (c *NewsController) CreateNews(ctx *gin.Context) {
	var news models.News
	if err := ctx.ShouldBindJSON(&news); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result := c.DB.Create(&news)
	if result.Error != nil {
		ctx.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	ctx.JSON(201, news)
}

func (c *NewsController) UpdateNews(ctx *gin.Context) {
	id := ctx.Param("id")
	var news models.News

	if err := c.DB.First(&news, id).Error; err != nil {
		ctx.JSON(404, gin.H{"error": "News not found"})
		return
	}

	if err := ctx.ShouldBindJSON(&news); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.DB.Save(&news)
	ctx.JSON(200, news)
}

func (c *NewsController) DeleteNews(ctx *gin.Context) {
	id := ctx.Param("id")

	result := c.DB.Delete(&models.News{}, id)
	if result.Error != nil {
		ctx.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "News deleted successfully"})
}

//
//func (c *NewsController) GetCategories(ctx *gin.Context) {
//	var categories []models.Category
//
//	result := c.DB.Find(&categories)
//	if result.Error != nil {
//		ctx.JSON(500, gin.H{"error": result.Error.Error()})
//		return
//	}
//
//	ctx.JSON(200, categories)
//}
//
//func (c *NewsController) GetSources(ctx *gin.Context) {
//	var sources []models.Source
//
//	result := c.DB.Find(&sources)
//	if result.Error != nil {
//		ctx.JSON(500, gin.H{"error": result.Error.Error()})
//		return
//	}
//
//	ctx.JSON(200, sources)
//}

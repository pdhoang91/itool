// 1. Package Declaration
package routes

// 2. Required Imports
import (
	"news-api/controllers" // Import local controller package

	"github.com/gin-gonic/gin" // Gin web framework
	"gorm.io/gorm"             // GORM database ORM
)

// 3. Route Setup Function
func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// 4. Initialize Controller
	newsController := &controllers.NewsController{DB: db}
	categoryController := &controllers.CategoryController{DB: db}
	sourceController := &controllers.SourceController{DB: db}

	// 5. Create API Group
	api := r.Group("/api")
	{
		// Category routes
		api.POST("/categories", categoryController.CreateCategory)
		api.GET("/categories", categoryController.GetCategories)

		// Source routes
		api.POST("/sources", sourceController.CreateSource)
		api.GET("/sources", sourceController.GetSources)

		// 6. Define Routes
		// GET /api/news - Get list of news with pagination
		api.GET("/news", newsController.GetNews)

		// GET /api/news/:id - Get single news detail
		// Example: /api/news/123
		api.GET("/news/:id", newsController.GetNewsDetail)

		// GET /api/news/source/:id - Get news by source with pagination
		// Example: /api/news/source/1
		api.GET("/news/source/:id", newsController.GetNewsBySource)

		// 7. Additional Routes Example
		// POST route for creating news
		api.POST("/news", newsController.CreateNews)

		// PUT route for updating news
		api.PUT("/news/:id", newsController.UpdateNews)

		// DELETE route for removing news
		api.DELETE("/news/:id", newsController.DeleteNews)

		//// GET route for categories
		//api.GET("/categories", newsController.GetCategories)
		//
		//// GET route for sources
		//api.GET("/sources", newsController.GetSources)
	}
}

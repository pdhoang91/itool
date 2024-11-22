// main.go
package main

import (
	"log"
	"news-api/config"
	"news-api/models"
	"news-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to database
	db := config.ConnectDB()

	// Auto migrate models
	db.AutoMigrate(&models.News{}, &models.Category{}, &models.Source{})

	// Initialize Gin router
	r := gin.Default()

	// corsConfig := cors.Config{
	// 	AllowOrigins:     []string{"http://202.92.6.77:3000", "http://localhost:3000", "https://insight.io.vn"},
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * 60 * 60, // 12 hours
	// }

	// r.Use(cors.New(corsConfig))

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Setup routes
	routes.SetupRoutes(r, db)

	// Start server
	if err := r.Run(":82"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

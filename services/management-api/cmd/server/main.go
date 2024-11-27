// cmd/main.go

package main

import (
	"log"

	"management-api/internal/config"
	"management-api/internal/handler"
	"management-api/internal/router"
	"management-api/internal/service"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	// Initialize service with both HTTP client and repository
	taskService := service.NewTaskService()

	// Initialize handler
	taskHandler := handler.NewTaskHandler(taskService)

	// Setup router
	r := router.SetupRouter(*taskHandler, cfg)

	// Run server
	if err := r.Run(cfg.Server.Port); err != nil {
		log.Fatalf("Could not run server: %v", err)
	}
}

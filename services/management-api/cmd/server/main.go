// cmd/main.go

package main

import (
	"log"

	"management-api/internal/config"
	"management-api/internal/handler"
	"management-api/internal/repository"
	"management-api/internal/router"
	"management-api/internal/service"

	"github.com/go-resty/resty/v2"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	// Initialize repository
	repo, err := repository.NewTaskRepository(cfg.Database)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer repo.Close()

	// Initialize service with both HTTP client and repository
	client := resty.New()
	taskService := service.NewTaskService(client, repo)

	// Initialize handler
	taskHandler := handler.NewTaskHandler(taskService)

	// Setup router
	r := router.SetupRouter(*taskHandler, cfg)

	// Run server
	if err := r.Run(cfg.Server.Port); err != nil {
		log.Fatalf("Could not run server: %v", err)
	}
}

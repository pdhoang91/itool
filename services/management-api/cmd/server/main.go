package main

import (
	"log"

	"management-api/internal/config"
	"management-api/internal/repository"
	"management-api/internal/router"
	"management-api/internal/service"
)

func main() {
	// Tải cấu hình
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	// Kết nối đến cơ sở dữ liệu
	repo, err := repository.NewTaskRepository(cfg.Database)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer repo.Close()

	// Khởi tạo service
	taskService := service.NewTaskService(repo)

	// Khởi tạo router
	r := router.SetupRouter(taskService, cfg)

	// Chạy server
	if err := r.Run(cfg.Server.Port); err != nil {
		log.Fatalf("Could not run server: %v", err)
	}
}

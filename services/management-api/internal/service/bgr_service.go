package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// HandleBackgroundRemoval loại bỏ nền ảnh sử dụng dịch vụ Background Removal
func (s *taskService) HandleBackgroundRemoval(ctx context.Context, imagePath string) (string, error) {
	log.Printf("Processing background removal for image: %s", imagePath)

	reqBody := map[string]string{"image_path": imagePath}

	resp, err := s.makeRequest(ctx, http.MethodPost, fmt.Sprintf("%s/remove-bg", BackgroundServiceURL), reqBody, nil)
	if err != nil {
		log.Printf("HandleBackgroundRemoval - makeRequest error: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		ProcessedImagePath string `json:"processed_image_path"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error decoding background removal response: %w", err)
	}

	log.Printf("Successfully processed image. Result: %s", result.ProcessedImagePath)
	return result.ProcessedImagePath, nil
}

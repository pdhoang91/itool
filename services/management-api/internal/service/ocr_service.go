package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// HandleOCR thực hiện OCR trên ảnh sử dụng dịch vụ OCR
func (s *taskService) HandleOCR(ctx context.Context, imagePath string) (map[string]string, error) {
	reqBody := map[string]string{"image_path": imagePath}

	resp, err := s.makeRequest(ctx, http.MethodPost, fmt.Sprintf("%s/ocr", OCRServiceURL), reqBody, nil)
	if err != nil {
		log.Printf("HandleOCR - makeRequest error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var ocrResp OCRResponse
	if err := json.NewDecoder(resp.Body).Decode(&ocrResp); err != nil {
		return nil, fmt.Errorf("error decoding OCR response: %w", err)
	}

	if ocrResp.Status != "success" {
		return nil, fmt.Errorf("OCR request failed with status: %s", ocrResp.Status)
	}

	return map[string]string{"text": ocrResp.Text}, nil
}

// OCRResponse đại diện cho phản hồi từ dịch vụ OCR
type OCRResponse struct {
	Text   string `json:"text"`
	Status string `json:"status"`
}

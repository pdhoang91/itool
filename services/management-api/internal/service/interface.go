package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	TTSServiceURL         = "http://text_to_speech_service:5001"
	VoiceToTextServiceURL = "http://voice_to_text_service:5002"
	BackgroundServiceURL  = "http://background_removal_service:5003"
	SpeechServiceURL      = "http://speech_recognition_service:5004"
	FaceServiceURL        = "http://face_recognition_service:5005"
	OCRServiceURL         = "http://ocr_service:5006"
	TranslationServiceURL = "http://translation_service:5007"
	DefaultTimeout        = 30 * time.Second
)

// TaskService định nghĩa giao diện cho các tác vụ
type TaskService interface {
	HandleTextToVoice(ctx context.Context, text, language, model string, speed, pitch, volume float64) (map[string]string, error)
	HandleVoiceToText(ctx context.Context, audioURL string) (map[string]string, error)
	GetAvailableLanguages(ctx context.Context) ([]string, error)
	GetAvailableVoices(ctx context.Context, language string) ([]Voice, error)
	HandleBackgroundRemoval(ctx context.Context, imagePath string) (string, error)
	HandleSpeechRecognition(ctx context.Context, audioURL string) (map[string]string, error)
	HandleFaceRecognition(ctx context.Context, imagePath string) (map[string]interface{}, error)
	HandleOCR(ctx context.Context, imagePath string) (map[string]string, error)
	HandleTranslation(ctx context.Context, text, destLang string) (map[string]string, error)
	UploadAudio(ctx context.Context, filePath string) (string, error)
	GetAvailableModels(ctx context.Context) (*ModelsResponse, error)
}

// taskService là triển khai cụ thể của TaskService
type taskService struct {
	httpClient *http.Client
}

// NewTaskService tạo một instance mới của TaskService với HTTP client mặc định
func NewTaskService() TaskService {
	return &taskService{
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

// makeRequest tạo và gửi một HTTP request với các tham số đã cho
func (s *taskService) makeRequest(ctx context.Context, method, url string, body interface{}, headers map[string]string) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	// Tạo một HTTP request mới với context
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	// Đặt header mặc định nếu có body
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Đặt thêm các header tùy chỉnh nếu có
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Thực thi HTTP request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing HTTP request: %w", err)
	}

	// Kiểm tra mã trạng thái HTTP
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("received non-2xx status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	return resp, nil
}

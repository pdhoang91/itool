package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

// Structs để ánh xạ các phản hồi JSON từ API

// TTSRequest đại diện cho payload yêu cầu Text-to-Speech
type TTSRequest struct {
	Text     string  `json:"text"`
	Speed    float64 `json:"speed"`
	Pitch    float64 `json:"pitch"`
	Model    string  `json:"model"`
	Language string  `json:"language"`
	Volume   float64 `json:"volume"`
}

// ModelDetails đại diện cho chi tiết một model TTS
type ModelDetails struct {
	Architecture string `json:"architecture"`
	Dataset      string `json:"dataset"`
	Language     string `json:"language"`
	ModelID      string `json:"model_id"`
	Type         string `json:"type"`
}

// ModelsResponse đại diện cho phản hồi từ API /models
type ModelsResponse struct {
	CurrentModel     string                    `json:"current_model"`
	ModelsByLanguage map[string][]ModelDetails `json:"models_by_language"`
	TotalModels      int                       `json:"total_models"`
}

// LanguageListResponse đại diện cho phản hồi từ API /models/languages
type LanguageListResponse struct {
	CurrentModel   string   `json:"current_model"`
	Languages      []string `json:"languages"`
	TotalLanguages int      `json:"total_languages"`
}

// TTSResponse đại diện cho phản hồi từ API /tts
type TTSResponse struct {
	AudioURL string `json:"audio_url"`
	FileName string `json:"filename"`
	Success  bool   `json:"success"`
}

// Voice đại diện cho một voice model
type Voice struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Gender   string `json:"gender,omitempty"`
	Model    string `json:"model"`
	Language string `json:"language"`
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

// HandleTextToVoice chuyển văn bản thành giọng nói sử dụng dịch vụ TTS
func (s *taskService) HandleTextToVoice(ctx context.Context, text, language, model string, speed, pitch, volume float64) (map[string]string, error) {
	ttsReq := TTSRequest{
		Text:     text,
		Speed:    speed,
		Pitch:    pitch,
		Model:    model,
		Language: language,
		Volume:   volume,
	}

	fmt.Println("ttsReq", ttsReq)

	resp, err := s.makeRequest(ctx, http.MethodPost, fmt.Sprintf("%s/tts", TTSServiceURL), ttsReq, nil)
	fmt.Println("resp", resp)
	if err != nil {
		log.Printf("HandleTextToVoice - makeRequest error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var ttsResp TTSResponse
	if err := json.NewDecoder(resp.Body).Decode(&ttsResp); err != nil {
		return nil, fmt.Errorf("error decoding TTS response: %w", err)
	}

	if !ttsResp.Success {
		return nil, fmt.Errorf("TTS request failed with status: %s", ttsResp.Success)
	}

	return map[string]string{"audio_url": "/mp3/" + ttsResp.FileName}, nil
}

// GetAvailableModels lấy danh sách các model TTS có sẵn từ dịch vụ TTS
func (s *taskService) GetAvailableModels(ctx context.Context) (*ModelsResponse, error) {
	resp, err := s.makeRequest(ctx, http.MethodGet, fmt.Sprintf("%s/models", TTSServiceURL), nil, nil)
	if err != nil {
		log.Printf("GetAvailableModels - makeRequest error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var modelsResp ModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&modelsResp); err != nil {
		return nil, fmt.Errorf("error decoding models response: %w", err)
	}

	return &modelsResp, nil
}

// GetAvailableLanguages lấy danh sách ngôn ngữ có sẵn từ dịch vụ TTS
func (s *taskService) GetAvailableLanguages(ctx context.Context) ([]string, error) {
	resp, err := s.makeRequest(ctx, http.MethodGet, fmt.Sprintf("%s/models/languages", TTSServiceURL), nil, nil)
	if err != nil {
		log.Printf("GetAvailableLanguages - makeRequest error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var langResp LanguageListResponse
	if err := json.NewDecoder(resp.Body).Decode(&langResp); err != nil {
		return nil, fmt.Errorf("error decoding language list response: %w", err)
	}

	return langResp.Languages, nil
}

// GetAvailableVoices lấy danh sách các voice có sẵn cho một ngôn ngữ cụ thể từ dịch vụ TTS
func (s *taskService) GetAvailableVoices(ctx context.Context, language string) ([]Voice, error) {
	modelsResp, err := s.GetAvailableModels(ctx)
	if err != nil {
		return nil, err
	}

	models, exists := modelsResp.ModelsByLanguage[language]
	if !exists {
		return nil, fmt.Errorf("no models found for language: %s", language)
	}

	var voices []Voice
	for _, model := range models {
		voice := Voice{
			ID:       model.ModelID,
			Name:     model.Architecture,
			Model:    model.ModelID,
			Language: model.Language,
		}
		voices = append(voices, voice)
	}

	return voices, nil
}

// HandleVoiceToText chuyển giọng nói thành văn bản sử dụng dịch vụ Voice-to-Text
func (s *taskService) HandleVoiceToText(ctx context.Context, audioURL string) (map[string]string, error) {
	reqBody := map[string]string{"audio_url": audioURL}

	resp, err := s.makeRequest(ctx, http.MethodPost, fmt.Sprintf("%s/convert", VoiceToTextServiceURL), reqBody, nil)
	if err != nil {
		log.Printf("HandleVoiceToText - makeRequest error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var vtResp VoiceToTextResponse
	if err := json.NewDecoder(resp.Body).Decode(&vtResp); err != nil {
		return nil, fmt.Errorf("error decoding Voice-to-Text response: %w", err)
	}

	if vtResp.Status != "success" {
		return nil, fmt.Errorf("Voice-to-Text request failed with status: %s", vtResp.Status)
	}

	return map[string]string{"transcript": vtResp.Transcript}, nil
}

// VoiceToTextResponse đại diện cho phản hồi từ dịch vụ Voice-to-Text
type VoiceToTextResponse struct {
	Transcript string `json:"transcript"`
	Status     string `json:"status"`
}

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

// HandleSpeechRecognition thực hiện nhận dạng giọng nói sử dụng dịch vụ Speech Recognition
func (s *taskService) HandleSpeechRecognition(ctx context.Context, audioURL string) (map[string]string, error) {
	reqBody := map[string]string{"audio_url": audioURL}

	resp, err := s.makeRequest(ctx, http.MethodPost, fmt.Sprintf("%s/recognize", SpeechServiceURL), reqBody, nil)
	if err != nil {
		log.Printf("HandleSpeechRecognition - makeRequest error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var srResp SpeechRecognitionResponse
	if err := json.NewDecoder(resp.Body).Decode(&srResp); err != nil {
		return nil, fmt.Errorf("error decoding speech recognition response: %w", err)
	}

	if srResp.Status != "success" {
		return nil, fmt.Errorf("Speech Recognition request failed with status: %s", srResp.Status)
	}

	return map[string]string{"transcript": srResp.Transcript}, nil
}

// SpeechRecognitionResponse đại diện cho phản hồi từ dịch vụ Speech Recognition
type SpeechRecognitionResponse struct {
	Transcript string `json:"transcript"`
	Status     string `json:"status"`
}

// HandleFaceRecognition thực hiện nhận dạng khuôn mặt sử dụng dịch vụ Face Recognition
func (s *taskService) HandleFaceRecognition(ctx context.Context, imagePath string) (map[string]interface{}, error) {
	reqBody := map[string]string{"image_path": imagePath}

	resp, err := s.makeRequest(ctx, http.MethodPost, fmt.Sprintf("%s/recognize-face", FaceServiceURL), reqBody, nil)
	if err != nil {
		log.Printf("HandleFaceRecognition - makeRequest error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var faceResp FaceRecognitionResponse
	if err := json.NewDecoder(resp.Body).Decode(&faceResp); err != nil {
		return nil, fmt.Errorf("error decoding face recognition response: %w", err)
	}

	return faceResp, nil
}

// FaceRecognitionResponse đại diện cho phản hồi từ dịch vụ Face Recognition
type FaceRecognitionResponse map[string]interface{}

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

// HandleTranslation thực hiện dịch văn bản sang ngôn ngữ đích sử dụng dịch vụ Translation
func (s *taskService) HandleTranslation(ctx context.Context, text, destLang string) (map[string]string, error) {
	reqBody := map[string]string{
		"text":      text,
		"dest_lang": destLang,
	}

	resp, err := s.makeRequest(ctx, http.MethodPost, fmt.Sprintf("%s/translate", TranslationServiceURL), reqBody, nil)
	if err != nil {
		log.Printf("HandleTranslation - makeRequest error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var translateResp TranslationResponse
	if err := json.NewDecoder(resp.Body).Decode(&translateResp); err != nil {
		return nil, fmt.Errorf("error decoding translation response: %w", err)
	}

	if translateResp.Status != "success" {
		return nil, fmt.Errorf("Translation request failed with status: %s", translateResp.Status)
	}

	return map[string]string{"translated_text": translateResp.TranslatedText}, nil
}

// TranslationResponse đại diện cho phản hồi từ dịch vụ Translation
type TranslationResponse struct {
	TranslatedText string `json:"translated_text"`
	Status         string `json:"status"`
}

// UploadAudio tải lên một file âm thanh và trả về URL của nó
func (s *taskService) UploadAudio(ctx context.Context, filePath string) (string, error) {
	// Triển khai logic tải lên thực tế ở đây
	// Hiện tại, trả về một URL giả định
	audioURL := fmt.Sprintf("http://localhost:81/uploads/audio/%s", filePath)
	return audioURL, nil
}

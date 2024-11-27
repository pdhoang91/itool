package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type TaskService interface {
	HandleTextToVoice(text, language, voice string, speed, pitch, volume float64) (map[string]string, error)
	HandleVoiceToText(audioURL string) (map[string]string, error)
	GetAvailableLanguages() ([]Language, error)
	GetAvailableVoices(language string) ([]Voice, error)
	HandleBackgroundRemoval(imagePath string) (string, error)
	HandleSpeechRecognition(audioURL string) (map[string]string, error)
	HandleFaceRecognition(imagePath string) (map[string]interface{}, error)
	HandleOCR(imagePath string) (map[string]string, error)
	HandleTranslation(text, destLang string) (map[string]string, error)
	UploadAudio(filePath string) (string, error)
}

type taskService struct {
	client *resty.Client
}

func NewTaskService(client *resty.Client) TaskService {
	return &taskService{
		client: client,
	}
}

const (
	TTS_SERVICE_URL = "http://localhost:5001"
)

type TTSResponse struct {
	AudioURL string `json:"audio_url"`
}

type Language struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Voice struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Model    string `json:"model"`
	Language string `json:"language"`
}

func (s *taskService) HandleTextToVoice(text, language, voice string, speed, pitch, volume float64) (map[string]string, error) {
	log.Printf("Starting TTS conversion: text=%s, language=%s, voice=%s", text, language, voice)

	// Create task record
	// task := &domain.Task{
	// 	ServiceName: "text-to-voice",
	// 	Status:      "processing",
	// 	InputData:   json.RawMessage(fmt.Sprintf(`{"text":"%s","language":"%s","voice":"%s","speed":%f,"pitch":%f,"volume":%f}`, text, language, voice, speed, pitch, volume)),
	// 	CreatedAt:   time.Now(),
	// 	UpdatedAt:   time.Now(),
	// }

	// Save initial task
	// taskID, err := s.repo.CreateTask(task)
	// if err != nil {
	// 	log.Printf("Failed to create task record: %v", err)
	// 	return nil, fmt.Errorf("failed to create task: %v", err)
	// }

	// Prepare request to TTS service
	ttsReq := map[string]interface{}{
		"text":     text,
		"language": language,
		"voice":    voice,
		"speed":    speed,
		"pitch":    pitch,
		"volume":   volume,
	}

	var ttsResp TTSResponse
	resp, err := s.client.R().
		SetBody(ttsReq).
		SetResult(&ttsResp).
		Post(fmt.Sprintf("%s/tts", TTS_SERVICE_URL))

	if err != nil || resp.IsError() {
		log.Printf("TTS service error: %v, response: %s", err, resp.String())
		// Update task status to failed
		//s.updateTaskStatus(taskID, "failed", nil)
		return nil, fmt.Errorf("TTS service error: %v", err)
	}

	// Update task status to completed
	outputData := map[string]string{"audio_url": ttsResp.AudioURL}
	// outputJSON, _ := json.Marshal(outputData)
	// err = s.updateTaskStatus(taskID, "completed", outputJSON)
	// if err != nil {
	// 	log.Printf("Failed to update task status: %v", err)
	// }

	return outputData, nil
}

func (s *taskService) GetAvailableLanguages() ([]Language, error) {
	var languages []Language

	resp, err := s.client.R().
		SetResult(&languages).
		Get(fmt.Sprintf("%s/languages", TTS_SERVICE_URL))

	if err != nil || resp.IsError() {
		log.Printf("Failed to fetch languages: %v", err)
		return nil, fmt.Errorf("failed to fetch languages: %v", err)
	}

	return languages, nil
}

func (s *taskService) GetAvailableVoices(language string) ([]Voice, error) {
	var voices []Voice

	resp, err := s.client.R().
		SetResult(&voices).
		Get(fmt.Sprintf("%s/voices/%s", TTS_SERVICE_URL, language))

	if err != nil || resp.IsError() {
		log.Printf("Failed to fetch voices for language %s: %v", language, err)
		return nil, fmt.Errorf("failed to fetch voices: %v", err)
	}

	return voices, nil
}

// HandleVoiceToText xử lý dịch vụ Voice-to-Text
func (s *taskService) HandleVoiceToText(audioURL string) (map[string]string, error) {
	resp, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{"audio_url": audioURL}).
		Post("http://voice_to_text_service:5002/convert")
	if err != nil || resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to call Voice-to-Text service")
	}

	var vtsResp map[string]string
	if err := json.Unmarshal(resp.Body(), &vtsResp); err != nil {
		return nil, fmt.Errorf("failed to parse Voice-to-Text response")
	}

	return vtsResp, nil
}

func (s *taskService) HandleBackgroundRemoval(imagePath string) (string, error) {
	log.Printf("HandleBackgroundRemoval: Received request with image path '%s'", imagePath)

	resp, err := s.client.R().
		SetFile("image", imagePath).
		Post("http://background_removal_service:5003/remove-bg")
	if err != nil {
		log.Printf("HandleBackgroundRemoval: Failed to call Background Removal service. Error: %v", err)
		return "", fmt.Errorf("failed to call Background Removal service")
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("HandleBackgroundRemoval: Service returned non-200 status. StatusCode: %d, Response: %s", resp.StatusCode(), resp.String())
		return "", fmt.Errorf("failed to call Background Removal service with StatusCode: %d", resp.StatusCode())
	}

	var brResp struct {
		ProcessedImagePath string `json:"processed_image_path"`
	}
	if err := json.Unmarshal(resp.Body(), &brResp); err != nil {
		log.Printf("HandleBackgroundRemoval: Failed to parse service response. Error: %v, Raw Response: %s", err, resp.String())
		return "", fmt.Errorf("failed to parse Background Removal response")
	}

	log.Printf("HandleBackgroundRemoval: Successfully processed background removal for image '%s'", imagePath)
	return brResp.ProcessedImagePath, nil
}

// HandleSpeechRecognition xử lý dịch vụ Speech Recognition
func (s *taskService) HandleSpeechRecognition(audioURL string) (map[string]string, error) {
	resp, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{"audio_url": audioURL}).
		Post("http://speech_recognition_service:5004/recognize")
	if err != nil || resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to call Speech Recognition service")
	}

	var srResp map[string]string
	if err := json.Unmarshal(resp.Body(), &srResp); err != nil {
		return nil, fmt.Errorf("failed to parse Speech Recognition response")
	}

	return srResp, nil
}

// HandleFaceRecognition xử lý dịch vụ Face Recognition
func (s *taskService) HandleFaceRecognition(imagePath string) (map[string]interface{}, error) {
	resp, err := s.client.R().
		SetFile("image", imagePath).
		Post("http://face_recognition_service:5005/recognize-face")
	if err != nil || resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to call Face Recognition service")
	}

	var frResp map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &frResp); err != nil {
		return nil, fmt.Errorf("failed to parse Face Recognition response")
	}

	return frResp, nil
}

// HandleOCR xử lý dịch vụ OCR
func (s *taskService) HandleOCR(imagePath string) (map[string]string, error) {
	resp, err := s.client.R().
		SetFile("image", imagePath).
		Post("http://ocr_service:5006/ocr")
	if err != nil || resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to call OCR service")
	}

	fmt.Println("resp", resp)

	var ocrResp map[string]string
	if err := json.Unmarshal(resp.Body(), &ocrResp); err != nil {
		return nil, fmt.Errorf("failed to parse OCR response")
	}
	fmt.Println("ocrResp", ocrResp)

	return ocrResp, nil
}

// HandleTranslation xử lý dịch vụ Translation
func (s *taskService) HandleTranslation(text, destLang string) (map[string]string, error) {
	resp, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{"text": text, "dest_lang": destLang}).
		Post("http://translation_service:5007/translate")
	if err != nil || resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to call Translation service")
	}

	var trResp map[string]string
	if err := json.Unmarshal(resp.Body(), &trResp); err != nil {
		return nil, fmt.Errorf("failed to parse Translation response")
	}

	return trResp, nil
}

// UploadAudio xử lý tải lên file audio
func (s *taskService) UploadAudio(filePath string) (string, error) {
	// Ở đây bạn có thể triển khai việc upload lên S3 hoặc dịch vụ lưu trữ khác.
	// Dưới đây là ví dụ trả về đường dẫn tạm thời.

	audioURL := fmt.Sprintf("http://localhost:81/uploads/audio/%s", filePath)
	return audioURL, nil
}

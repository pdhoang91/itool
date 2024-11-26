package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"management-api/internal/domain"
	"net/http"
)

func (s *taskService) HandleTextToVoice(text, language, voice string, speed, pitch, volume float64) (map[string]string, error) {
	log.Printf("HandleTextToVoice: Starting conversion with text '%s', language '%s', voice '%s'",
		text, language, voice)

	// Call TTS service
	resp, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"text":     text,
			"language": language,
			"voice":    voice,
			"speed":    speed,
			"pitch":    pitch,
			"volume":   volume,
		}).
		Post("http://text_to_voice_service:5001/convert")

	if err != nil || resp.StatusCode() != 200 {
		log.Printf("HandleTextToVoice: Service error - %v, Status: %d", err, resp.StatusCode())
		return nil, fmt.Errorf("failed to call Text-to-Voice service")
	}

	var ttsResp map[string]string
	if err := json.Unmarshal(resp.Body(), &ttsResp); err != nil {
		log.Printf("HandleTextToVoice: Response parsing error - %v", err)
		return nil, fmt.Errorf("failed to parse Text-to-Voice response")
	}

	// Create task record
	task := &domain.Task{
		ServiceName: "text-to-voice",
		Status:      "completed",
		InputData: json.RawMessage(fmt.Sprintf(`{"text":"%s","language":"%s","voice":"%s","speed":%f,"pitch":%f,"volume":%f}`,
			text, language, voice, speed, pitch, volume)),
		OutputData: json.RawMessage(fmt.Sprintf(`{"audio_url":"%s"}`, ttsResp["audio_url"])),
	}

	// Add CreateTask method to repository interface and implementation if not exists
	if err := s.repo.CreateTask(context.Background(), task); err != nil {
		log.Printf("HandleTextToVoice: Failed to create task record - %v", err)
		// Continue even if task record creation fails
	}

	log.Printf("HandleTextToVoice: Successfully converted text to voice")
	return ttsResp, nil
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

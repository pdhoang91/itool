package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// // HandleTextToVoice xử lý dịch vụ Text-to-Voice
//
//	func (s *taskService) HandleTextToVoice(text, language string) (map[string]string, error) {
//		if language == "" {
//			language = "en"
//		}
//
//		resp, err := s.client.R().
//			SetHeader("Content-Type", "application/json").
//			SetBody(map[string]string{"text": text, "language": language}).
//			Post("http://text_to_voice_service:5001/convert")
//		if err != nil || resp.StatusCode() != 200 {
//			return nil, fmt.Errorf("failed to call Text-to-Voice service")
//		}
//
//		var ttsResp map[string]string
//		if err := json.Unmarshal(resp.Body(), &ttsResp); err != nil {
//			return nil, fmt.Errorf("failed to parse Text-to-Voice response")
//		}
//
//		return ttsResp, nil
//	}
//
// HandleTextToVoice xử lý dịch vụ Text-to-Voice
func (s *taskService) HandleTextToVoice(text, language string) (map[string]string, error) {
	log.Printf("HandleTextToVoice: Received request with text '%s' and language '%s'", text, language)

	if language == "" {
		language = "en"
		log.Printf("HandleTextToVoice: Defaulting language to '%s'", language)
	}

	resp, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{"text": text, "language": language}).
		Post("http://text_to_voice_service:5001/convert")
	if err != nil || resp.StatusCode() != 200 {
		log.Printf("HandleTextToVoice: Error calling service. StatusCode: %d, Error: %v", resp.StatusCode(), err)
		return nil, fmt.Errorf("failed to call Text-to-Voice service")
	}

	var ttsResp map[string]string
	if err := json.Unmarshal(resp.Body(), &ttsResp); err != nil {
		log.Printf("HandleTextToVoice: Error parsing response. Error: %v", err)
		return nil, fmt.Errorf("failed to parse Text-to-Voice response")
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

	var ocrResp map[string]string
	if err := json.Unmarshal(resp.Body(), &ocrResp); err != nil {
		return nil, fmt.Errorf("failed to parse OCR response")
	}

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

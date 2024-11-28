package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

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

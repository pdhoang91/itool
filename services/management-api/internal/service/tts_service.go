package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

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

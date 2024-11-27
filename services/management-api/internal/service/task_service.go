package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	TTSServiceURL         = "http://localhost:5001"
	VoiceToTextServiceURL = "http://voice_to_text_service:5002"
	BackgroundServiceURL  = "http://background_removal_service:5003"
	SpeechServiceURL      = "http://speech_recognition_service:5004"
	FaceServiceURL        = "http://face_recognition_service:5005"
	OCRServiceURL         = "http://ocr_service:5006"
	TranslationServiceURL = "http://translation_service:5007"
	DefaultTimeout        = 30 * time.Second
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
	httpClient *http.Client
}

func NewTaskService() TaskService {
	return &taskService{
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

type TTSRequest struct {
	Text     string  `json:"text"`
	Speed    float64 `json:"speed"`
	Pitch    float64 `json:"pitch"`
	Model    string  `json:"model"`
	Language string  `json:"language"`
	Volume   float64 `json:"volume"`
}

type Language struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type LanguageResponse struct {
	Code      string     `json:"code"`
	Name      string     `json:"name"`
	Languages []Language `json:"language"`
}

type ModelsResponse struct {
	Code             string      `json:"code"`
	Name             string      `json:"name"`
	Languages        []Language  `json:"language"`
	ModelsByLanguage interface{} `json:"modelsByLanguage"`
}

type Voice struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Model    string `json:"model"`
	Language string `json:"language"`
}

func (s *taskService) makeRequest(method, url string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("service returned status: %d", resp.StatusCode)
	}

	return resp, nil
}

func (s *taskService) HandleTextToVoice(text, language, voice string, speed, pitch, volume float64) (map[string]string, error) {
	ttsReq := TTSRequest{
		Text:     text,
		Speed:    speed,
		Pitch:    pitch,
		Model:    voice,
		Language: language,
		Volume:   volume,
	}

	resp, err := s.makeRequest(http.MethodPost, fmt.Sprintf("%s/tts", TTSServiceURL), ttsReq)
	fmt.Println("resp", resp)
	if err != nil {
		fmt.Println("err", err)
		return nil, err
	}
	defer resp.Body.Close()

	audioURL := fmt.Sprintf("http://localhost:81/uploads/audio/output_%d.wav", time.Now().Unix())
	return map[string]string{"audio_url": audioURL}, nil
}

func (s *taskService) GetAvailableLanguages() ([]Language, error) {
	resp, err := s.makeRequest(http.MethodGet, fmt.Sprintf("%s/models/languages", TTSServiceURL), nil)
	fmt.Println("resp", resp)
	if err != nil {
		fmt.Println("err", err)
		return nil, err
	}
	defer resp.Body.Close()

	var langResp LanguageResponse
	if err := json.NewDecoder(resp.Body).Decode(&langResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	languages := make([]Language, len(langResp.Languages))
	for i, lang := range langResp.Languages {
		languages[i] = Language{
			Code: lang.Code,
			Name: lang.Name,
		}
	}

	return languages, nil
}

func (s *taskService) GetAvailableVoices(language string) ([]Voice, error) {
	resp, err := s.makeRequest(http.MethodGet, fmt.Sprintf("%s/models", TTSServiceURL), nil)
	fmt.Println("resp", resp)
	if err != nil {
		fmt.Println("err", err)
		return nil, err
	}
	defer resp.Body.Close()

	var modelsResp ModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&modelsResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	var voices []Voice
	// if models, ok := modelsResp.ModelsByLanguage[language]; ok {
	// 	voices = make([]Voice, len(models))
	// 	for i, model := range models {
	// 		voices[i] = Voice{
	// 			ID:       model.ModelID,
	// 			Name:     model.Architecture,
	// 			Language: model.Language,
	// 			Model:    model.ModelID,
	// 		}
	// 	}
	// }

	return voices, nil
}

func (s *taskService) HandleVoiceToText(audioURL string) (map[string]string, error) {
	resp, err := s.makeRequest(http.MethodPost, fmt.Sprintf("%s/convert", VoiceToTextServiceURL),
		map[string]string{"audio_url": audioURL})

	fmt.Println("resp", resp)
	if err != nil {
		fmt.Println("err", err)
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return result, nil
}

func (s *taskService) HandleBackgroundRemoval(imagePath string) (string, error) {
	log.Printf("Processing background removal for image: %s", imagePath)

	resp, err := s.makeRequest(http.MethodPost, fmt.Sprintf("%s/remove-bg", BackgroundServiceURL),
		map[string]string{"image_path": imagePath})
	if err != nil {
		log.Printf("Background removal failed: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	fmt.Println("resp", resp)

	var result struct {
		ProcessedImagePath string `json:"processed_image_path"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	log.Printf("Successfully processed image. Result: %s", result.ProcessedImagePath)
	return result.ProcessedImagePath, nil
}

func (s *taskService) HandleSpeechRecognition(audioURL string) (map[string]string, error) {
	resp, err := s.makeRequest(http.MethodPost, fmt.Sprintf("%s/recognize", SpeechServiceURL),
		map[string]string{"audio_url": audioURL})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println("resp", resp)
	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return result, nil
}

func (s *taskService) HandleFaceRecognition(imagePath string) (map[string]interface{}, error) {
	resp, err := s.makeRequest(http.MethodPost, fmt.Sprintf("%s/recognize-face", FaceServiceURL),
		map[string]string{"image_path": imagePath})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return result, nil
}

func (s *taskService) HandleOCR(imagePath string) (map[string]string, error) {
	resp, err := s.makeRequest(http.MethodPost, fmt.Sprintf("%s/ocr", OCRServiceURL),
		map[string]string{"image_path": imagePath})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return result, nil
}

func (s *taskService) HandleTranslation(text, destLang string) (map[string]string, error) {
	resp, err := s.makeRequest(http.MethodPost, fmt.Sprintf("%s/translate", TranslationServiceURL),
		map[string]string{
			"text":      text,
			"dest_lang": destLang,
		})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return result, nil
}

func (s *taskService) UploadAudio(filePath string) (string, error) {
	audioURL := fmt.Sprintf("http://localhost:81/uploads/audio/%s", filePath)
	return audioURL, nil
}

package handler

import (
	"log"
	"net/http"

	"management-api/internal/service"
	"management-api/pkg/utils"

	"github.com/gin-gonic/gin"
)

// TaskHandler định nghĩa handler cho các endpoint liên quan đến tác vụ
type TaskHandler struct {
	service service.TaskService
}

// NewTaskHandler tạo một instance mới của TaskHandler với TaskService đã cung cấp
func NewTaskHandler(service service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

// GetAvailableLanguages xử lý endpoint GET /languages
func (h *TaskHandler) GetAvailableLanguages(c *gin.Context) {
	log.Println("GetAvailableLanguages: Received request")

	// Lấy context từ request
	ctx := c.Request.Context()

	// Gọi service với context
	languages, err := h.service.GetAvailableLanguages(ctx)
	if err != nil {
		log.Printf("GetAvailableLanguages: Error fetching languages - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch languages"})
		return
	}

	log.Printf("GetAvailableLanguages: Successfully retrieved %d languages", len(languages))
	c.JSON(http.StatusOK, gin.H{
		"languages":       languages,
		"total_languages": len(languages),
	})
}

// GetAvailableVoices xử lý endpoint GET /voices/:language
func (h *TaskHandler) GetAvailableVoices(c *gin.Context) {
	language := c.Param("language")
	log.Printf("GetAvailableVoices: Received request for language: %s", language)

	// Lấy context từ request
	ctx := c.Request.Context()

	// Gọi service với context
	voices, err := h.service.GetAvailableVoices(ctx, language)
	if err != nil {
		log.Printf("GetAvailableVoices: Error fetching voices - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch voices"})
		return
	}

	log.Printf("GetAvailableVoices: Successfully retrieved %d voices for language %s", len(voices), language)
	c.JSON(http.StatusOK, gin.H{
		"voices": voices,
	})
}

// HandleTextToVoice xử lý endpoint POST /tts
func (h *TaskHandler) HandleTextToVoice(c *gin.Context) {
	log.Println("HandleTextToVoice: Received request")

	var req struct {
		Text     string  `json:"text" binding:"required"`
		Language string  `json:"language"`
		Model    string  `json:"model"`
		Speed    float64 `json:"speed"`
		Pitch    float64 `json:"pitch"`
		Volume   float64 `json:"volume"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("HandleTextToVoice: Invalid input - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Đặt giá trị mặc định nếu cần thiết
	if req.Language == "" {
		req.Language = "en"
	}
	if req.Speed == 0 {
		req.Speed = 1.0
	}
	if req.Pitch == 0 {
		req.Pitch = 0.0
	}
	if req.Volume == 0 {
		req.Volume = 1.0
	}
	if req.Model == "" {
		req.Model = "tts_models/en/ljspeech/tacotron2-DDC"
	}

	log.Printf("HandleTextToVoice: Processing request - Text: %s, Language: %s, Model: %s",
		req.Text, req.Language, req.Model)

	// Lấy context từ request
	ctx := c.Request.Context()

	// Gọi service với context
	resp, err := h.service.HandleTextToVoice(
		ctx,
		req.Text,
		req.Language,
		req.Model,
		req.Speed,
		req.Pitch,
		req.Volume,
	)
	if err != nil {
		log.Printf("HandleTextToVoice: Service error - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("HandleTextToVoice: Success - Response: %v", resp)
	c.JSON(http.StatusOK, resp)
}

// HandleVoiceToText xử lý endpoint POST /vts
func (h *TaskHandler) HandleVoiceToText(c *gin.Context) {
	var req struct {
		AudioURL string `json:"audio_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("HandleVoiceToText: Invalid input - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Lấy context từ request
	ctx := c.Request.Context()

	// Gọi service với context
	resp, err := h.service.HandleVoiceToText(ctx, req.AudioURL)
	if err != nil {
		log.Printf("HandleVoiceToText: Service error - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("HandleVoiceToText: Success - Response: %v", resp)
	c.JSON(http.StatusOK, resp)
}

// HandleBackgroundRemoval xử lý endpoint POST /remove-bg
func (h *TaskHandler) HandleBackgroundRemoval(c *gin.Context) {
	log.Println("HandleBackgroundRemoval: Received request to remove background")

	// Lấy file từ yêu cầu
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		log.Printf("HandleBackgroundRemoval: No image file provided. Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided"})
		return
	}
	log.Printf("HandleBackgroundRemoval: Received file '%s' with size %d bytes", header.Filename, header.Size)

	// Lưu file tạm thời
	uploadPath := "./uploads/images/"
	filePath, err := utils.SaveUploadedFile(file, header, uploadPath)
	if err != nil {
		log.Printf("HandleBackgroundRemoval: Failed to save file '%s'. Error: %v", header.Filename, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save the file"})
		return
	}
	log.Printf("HandleBackgroundRemoval: File saved to temporary path '%s'", filePath)

	// Lấy context từ request
	ctx := c.Request.Context()

	// Gọi service xử lý background removal với context
	processedImagePath, err := h.service.HandleBackgroundRemoval(ctx, filePath)
	if err != nil {
		log.Printf("HandleBackgroundRemoval: Failed to process background removal for file '%s'. Error: %v", filePath, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("HandleBackgroundRemoval: Successfully processed background removal for file '%s'", filePath)

	// Trả về đường dẫn file đã xử lý
	c.JSON(http.StatusOK, gin.H{"processed_image_path": processedImagePath})
}

// HandleSpeechRecognition xử lý endpoint POST /speech-recognition
func (h *TaskHandler) HandleSpeechRecognition(c *gin.Context) {
	var req struct {
		AudioURL string `json:"audio_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("HandleSpeechRecognition: Invalid input - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Lấy context từ request
	ctx := c.Request.Context()

	// Gọi service với context
	resp, err := h.service.HandleSpeechRecognition(ctx, req.AudioURL)
	if err != nil {
		log.Printf("HandleSpeechRecognition: Service error - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("HandleSpeechRecognition: Success - Response: %v", resp)
	c.JSON(http.StatusOK, resp)
}

// HandleFaceRecognition xử lý endpoint POST /face-recognition
func (h *TaskHandler) HandleFaceRecognition(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		log.Printf("HandleFaceRecognition: No image file provided. Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided"})
		return
	}

	// Lưu file tạm thời
	uploadPath := "./uploads/images/"
	filePath, err := utils.SaveUploadedFile(file, header, uploadPath)
	if err != nil {
		log.Printf("HandleFaceRecognition: Failed to save file '%s'. Error: %v", header.Filename, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save the file"})
		return
	}
	log.Printf("HandleFaceRecognition: File saved to temporary path '%s'", filePath)

	// Lấy context từ request
	ctx := c.Request.Context()

	// Gọi service xử lý face recognition với context
	resp, err := h.service.HandleFaceRecognition(ctx, filePath)
	if err != nil {
		log.Printf("HandleFaceRecognition: Failed to process face recognition for file '%s'. Error: %v", filePath, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("HandleFaceRecognition: Successfully processed face recognition for file '%s'", filePath)

	c.JSON(http.StatusOK, resp)
}

// HandleOCR xử lý endpoint POST /ocr
func (h *TaskHandler) HandleOCR(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		log.Printf("HandleOCR: No image file provided. Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided"})
		return
	}

	// Lưu file tạm thời
	uploadPath := "./uploads/images/"
	filePath, err := utils.SaveUploadedFile(file, header, uploadPath)
	if err != nil {
		log.Printf("HandleOCR: Failed to save file '%s'. Error: %v", header.Filename, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save the file"})
		return
	}
	log.Printf("HandleOCR: File saved to temporary path '%s'", filePath)

	// Lấy context từ request
	ctx := c.Request.Context()

	// Gọi service xử lý OCR với context
	resp, err := h.service.HandleOCR(ctx, filePath)
	if err != nil {
		log.Printf("HandleOCR: Failed to process OCR for file '%s'. Error: %v", filePath, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("HandleOCR: Successfully processed OCR for file '%s'", filePath)
	c.JSON(http.StatusOK, resp)
}

// HandleTranslation xử lý endpoint POST /translate
func (h *TaskHandler) HandleTranslation(c *gin.Context) {
	var req struct {
		Text     string `json:"text" binding:"required"`
		DestLang string `json:"dest_lang" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("HandleTranslation: Invalid input - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	log.Printf("HandleTranslation: Processing translation - Text: %s, DestLang: %s", req.Text, req.DestLang)

	// Lấy context từ request
	ctx := c.Request.Context()

	// Gọi service với context
	resp, err := h.service.HandleTranslation(ctx, req.Text, req.DestLang)
	if err != nil {
		log.Printf("HandleTranslation: Service error - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("HandleTranslation: Success - Response: %v", resp)
	c.JSON(http.StatusOK, resp)
}

// UploadAudio xử lý endpoint POST /upload-audio
func (h *TaskHandler) UploadAudio(c *gin.Context) {
	file, header, err := c.Request.FormFile("audio")
	if err != nil {
		log.Printf("UploadAudio: No audio file provided. Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No audio file provided"})
		return
	}

	// Lưu file vào thư mục /uploads/audio (tạo thư mục này nếu chưa có)
	uploadPath := "./uploads/audio/"
	filePath, err := utils.SaveUploadedFile(file, header, uploadPath)
	if err != nil {
		log.Printf("UploadAudio: Failed to save audio file '%s'. Error: %v", header.Filename, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save the audio file"})
		return
	}
	log.Printf("UploadAudio: Audio file saved to '%s'", filePath)

	// Lấy context từ request
	ctx := c.Request.Context()

	// Gọi service để tải lên (nếu cần, ví dụ như lên S3)
	audioURL, err := h.service.UploadAudio(ctx, header.Filename)
	if err != nil {
		log.Printf("UploadAudio: Failed to upload audio file '%s'. Error: %v", header.Filename, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload audio"})
		return
	}

	log.Printf("UploadAudio: Successfully uploaded audio. URL: %s", audioURL)
	c.JSON(http.StatusOK, gin.H{"audio_url": audioURL})
}

package handler

import (
	"log"
	"net/http"

	"management-api/internal/service"
	"management-api/pkg/utils"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	service service.TaskService
}

func NewTaskHandler(service service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) GetAvailableLanguages(c *gin.Context) {
	log.Println("GetAvailableLanguages: Received request")

	languages, err := h.service.GetAvailableLanguages()
	if err != nil {
		log.Printf("GetAvailableLanguages: Error fetching languages - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch languages"})
		return
	}

	log.Printf("GetAvailableLanguages: Successfully retrieved %d languages", len(languages))
	c.JSON(http.StatusOK, languages)
}

func (h *TaskHandler) GetAvailableVoices(c *gin.Context) {
	language := c.Param("language")
	log.Printf("GetAvailableVoices: Received request for language: %s", language)

	voices, err := h.service.GetAvailableVoices(language)
	if err != nil {
		log.Printf("GetAvailableVoices: Error fetching voices - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch voices"})
		return
	}

	log.Printf("GetAvailableVoices: Successfully retrieved %d voices for language %s", len(voices), language)
	c.JSON(http.StatusOK, voices)
}

func (h *TaskHandler) HandleTextToVoice(c *gin.Context) {
	log.Println("HandleTextToVoice: Received request")

	var req struct {
		Text     string  `json:"text" binding:"required"`
		Language string  `json:"language"`
		Voice    string  `json:"voice"`
		Speed    float64 `json:"speed"`
		Pitch    float64 `json:"pitch"`
		Volume   float64 `json:"volume"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("HandleTextToVoice: Invalid input - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Set default values
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

	log.Printf("HandleTextToVoice: Processing request - Text: %s, Language: %s, Voice: %s",
		req.Text, req.Language, req.Voice)

	resp, err := h.service.HandleTextToVoice(
		req.Text,
		req.Language,
		req.Voice,
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

// HandleVoiceToText xử lý endpoint /vts
func (h *TaskHandler) HandleVoiceToText(c *gin.Context) {
	var req struct {
		AudioURL string `json:"audio_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	resp, err := h.service.HandleVoiceToText(req.AudioURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// HandleBackgroundRemoval xử lý endpoint /remove-bg
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

	// Gọi service xử lý background removal
	processedImagePath, err := h.service.HandleBackgroundRemoval(filePath)
	if err != nil {
		log.Printf("HandleBackgroundRemoval: Failed to process background removal for file '%s'. Error: %v", filePath, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("HandleBackgroundRemoval: Successfully processed background removal for file '%s'", filePath)

	// Trả về đường dẫn file đã xử lý
	c.JSON(http.StatusOK, gin.H{"processed_image_path": processedImagePath})
}

// HandleSpeechRecognition xử lý endpoint /speech-recognition
func (h *TaskHandler) HandleSpeechRecognition(c *gin.Context) {
	var req struct {
		AudioURL string `json:"audio_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	resp, err := h.service.HandleSpeechRecognition(req.AudioURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// HandleFaceRecognition xử lý endpoint /face-recognition
func (h *TaskHandler) HandleFaceRecognition(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided"})
		return
	}

	// Lưu file tạm thời
	uploadPath := "./uploads/images/"
	filePath, err := utils.SaveUploadedFile(file, header, uploadPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save the file"})
		return
	}

	resp, err := h.service.HandleFaceRecognition(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// HandleOCR xử lý endpoint /ocr
func (h *TaskHandler) HandleOCR(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided"})
		return
	}

	// Lưu file tạm thời
	uploadPath := "./uploads/images/"
	filePath, err := utils.SaveUploadedFile(file, header, uploadPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save the file"})
		return
	}

	resp, err := h.service.HandleOCR(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// HandleTranslation xử lý endpoint /translate
func (h *TaskHandler) HandleTranslation(c *gin.Context) {
	var req struct {
		Text     string `json:"text" binding:"required"`
		DestLang string `json:"dest_lang" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	resp, err := h.service.HandleTranslation(req.Text, req.DestLang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UploadAudio xử lý endpoint /upload-audio
func (h *TaskHandler) UploadAudio(c *gin.Context) {
	file, header, err := c.Request.FormFile("audio")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No audio file provided"})
		return
	}

	// Lưu file vào thư mục /uploads/audio (tạo thư mục này nếu chưa có)
	uploadPath := "./uploads/audio/"
	filePath, err := utils.SaveUploadedFile(file, header, uploadPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save the file"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"audio_url": filePath})

	// Có thể triển khai upload lên S3 hoặc dịch vụ lưu trữ khác tại đây
	// Ví dụ trả về đường dẫn tạm thời
	//audioURL, err := h.service.UploadAudio(header.Filename)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload audio"})
	//	return
	//}

	//c.JSON(http.StatusOK, gin.H{"audio_url": audioURL})
}

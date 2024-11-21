package handler

import (
	"log"
	"net/http"
	"strconv"

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

// GetTaskStatus lấy trạng thái của một task
func (h *TaskHandler) GetTaskStatus(c *gin.Context) {
	idParam := c.Param("id")
	log.Printf("GetTaskStatus: Received request with ID %s", idParam)

	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Printf("GetTaskStatus: Invalid task ID %s. Error: %v", idParam, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	task, err := h.service.GetTaskStatus(id)
	if err != nil {
		log.Printf("GetTaskStatus: Task with ID %d not found. Error: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	log.Printf("GetTaskStatus: Successfully retrieved task ID %d", id)
	c.JSON(http.StatusOK, task)
}

// GetAllTasks lấy danh sách tất cả các task
func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	log.Println("GetAllTasks: Received request for all tasks")

	tasks, err := h.service.GetAllTasks()
	if err != nil {
		log.Printf("GetAllTasks: Failed to retrieve tasks. Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
		return
	}

	log.Printf("GetAllTasks: Successfully retrieved %d tasks", len(tasks))
	c.JSON(http.StatusOK, tasks)
}

// HandleTextToVoice xử lý endpoint /tts
func (h *TaskHandler) HandleTextToVoice(c *gin.Context) {
	var req struct {
		Text     string `json:"text" binding:"required"`
		Language string `json:"language"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	resp, err := h.service.HandleTextToVoice(req.Text, req.Language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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

//// HandleBackgroundRemoval xử lý endpoint /remove-bg
//func (h *TaskHandler) HandleBackgroundRemoval(c *gin.Context) {
//	file, header, err := c.Request.FormFile("image")
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided"})
//		return
//	}
//
//	// Lưu file tạm thời
//	uploadPath := "./uploads/images/"
//	filePath, err := utils.SaveUploadedFile(file, header, uploadPath)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save the file"})
//		return
//	}
//
//	resp, err := h.service.HandleBackgroundRemoval(filePath)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, resp)
//}

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

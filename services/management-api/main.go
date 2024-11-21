package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Task struct {
	ID          int             `json:"id"`
	ServiceName string          `json:"service_name"`
	Status      string          `json:"status"`
	InputData   json.RawMessage `json:"input_data"`
	OutputData  json.RawMessage `json:"output_data"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

var dbPool *pgxpool.Pool
var client *resty.Client

func main() {
	// os.Setenv("DB_HOST", "localhost")
	// os.Setenv("DB_PORT", "5433")
	// os.Setenv("DB_USER", "admin")
	// os.Setenv("DB_PASSWORD", "password")
	// os.Setenv("DB_NAME", "ai_tools")

	var err error
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	dbPool, err = pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		panic(err)
	}
	defer dbPool.Close()

	client = resty.New()

	r := gin.Default()

	config := cors.Config{
		//AllowOrigins:     []string{allowOrigins, "http://localhost:3000", "http://202.92.6.77:3000"}, // Thay đổi tùy vào frontend
		AllowOrigins:     []string{"http://202.92.6.77:3000", "http://localhost:3000", "https://insight.io.vn"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	r.Use(cors.New(config))

	r.Static("/uploads", "./uploads")

	// Các endpoint tương ứng với từng service
	r.POST("/tts", handleTextToVoice)
	r.POST("/vts", handleVoiceToText)
	r.POST("/remove-bg", handleBackgroundRemoval)
	r.POST("/speech-recognition", handleSpeechRecognition)
	r.POST("/face-recognition", handleFaceRecognition)
	r.POST("/ocr", handleOCR)
	r.POST("/translate", handleTranslation)
	r.POST("/upload-audio", uploadAudio)

	r.GET("/tasks/:id", getTaskStatus)
	r.GET("/tasks", getAllTasks)

	r.Run(":81")
}

// Thêm hàm uploadAudio
func uploadAudio(c *gin.Context) {
	file, header, err := c.Request.FormFile("audio")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No audio file provided"})
		return
	}

	// Lưu file vào thư mục /uploads/audio (tạo thư mục này nếu chưa có)
	uploadPath := "./uploads/audio/"
	os.MkdirAll(uploadPath, os.ModePerm)
	filePath := filepath.Join(uploadPath, header.Filename)

	out, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create the file"})
		return
	}
	defer out.Close()

	_, err = ioutil.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read the file"})
		return
	}

	// Bạn có thể sử dụng thư viện như aws-sdk-go để upload file lên S3 hoặc một dịch vụ lưu trữ file khác
	// Ở đây, tôi sẽ trả về đường dẫn tạm thời
	audioURL := fmt.Sprintf("http://localhost:81/uploads/audio/%s", header.Filename)

	c.JSON(http.StatusOK, gin.H{"audio_url": audioURL})
}

// Thêm hàm getAllTasks
func getAllTasks(c *gin.Context) {
	rows, err := dbPool.Query(context.Background(), "SELECT id, service_name, status, input_data, output_data, created_at, updated_at FROM tasks ORDER BY created_at DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
		return
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.ServiceName, &task.Status, &task.InputData, &task.OutputData, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database scan error"})
			return
		}
		tasks = append(tasks, task)
	}

	c.JSON(http.StatusOK, tasks)
}

func getTaskStatus(c *gin.Context) {
	taskID := c.Param("id")
	var task Task
	err := dbPool.QueryRow(context.Background(),
		"SELECT id, service_name, status, input_data, output_data, created_at, updated_at FROM tasks WHERE id=$1",
		taskID,
	).Scan(&task.ID, &task.ServiceName, &task.Status, &task.InputData, &task.OutputData, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func handleTextToVoice(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	text, textExists := req["text"]
	if !textExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'text' field"})
		return
	}

	// Default language is "en" if not provided
	language := req["language"]
	if language == "" {
		language = "en"
	}

	// Call Text-to-Voice service with text and language
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{"text": text, "language": language}).
		Post("http://text_to_voice_service:5001/convert")
	if err != nil || resp.StatusCode() != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call Text-to-Voice service"})
		return
	}

	var ttsResp map[string]string
	if err := json.Unmarshal(resp.Body(), &ttsResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse Text-to-Voice response"})
		return
	}

	c.JSON(http.StatusOK, ttsResp)
}

func handleVoiceToText(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	audioURL, exists := req["audio_url"]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'audio_url' field"})
		return
	}

	// Gọi service Voice-to-Text
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{"audio_url": audioURL}).
		Post("http://localhost:5002/convert")
	//Post("http://voice_to_text_service:5002/convert")
	if err != nil || resp.StatusCode() != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call Voice-to-Text service"})
		return
	}

	var vtsResp map[string]string
	json.Unmarshal(resp.Body(), &vtsResp)

	c.JSON(http.StatusOK, vtsResp)
}

func handleBackgroundRemoval(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided"})
		return
	}

	// Tạo request đến service Background Removal
	resp, err := client.R().
		SetFileReader("image", header.Filename, file).
		Post("http://background_removal_service:5003/remove-bg")
	if err != nil || resp.StatusCode() != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call Background Removal service"})
		return
	}

	var brResp map[string]string
	json.Unmarshal(resp.Body(), &brResp)

	c.JSON(http.StatusOK, brResp)
}

func handleSpeechRecognition(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	audioURL, exists := req["audio_url"]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'audio_url' field"})
		return
	}

	// Gọi service Speech Recognition
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{"audio_url": audioURL}).
		Post("http://speech_recognition_service:5004/recognize")
	if err != nil || resp.StatusCode() != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call Speech Recognition service"})
		return
	}

	var srResp map[string]string
	json.Unmarshal(resp.Body(), &srResp)

	c.JSON(http.StatusOK, srResp)
}

func handleFaceRecognition(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided"})
		return
	}

	// Tạo request đến service Face Recognition
	resp, err := client.R().
		SetFileReader("image", header.Filename, file).
		Post("http://face_recognition:5005/recognize-face")
	if err != nil || resp.StatusCode() != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call Face Recognition service"})
		return
	}

	var frResp map[string]interface{}
	json.Unmarshal(resp.Body(), &frResp)

	c.JSON(http.StatusOK, frResp)
}

func handleOCR(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided"})
		return
	}

	// Tạo request đến service OCR
	resp, err := client.R().
		SetFileReader("image", header.Filename, file).
		Post("http://ocr_service:5006/ocr")
	if err != nil || resp.StatusCode() != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call OCR service"})
		return
	}

	var ocrResp map[string]string
	json.Unmarshal(resp.Body(), &ocrResp)

	c.JSON(http.StatusOK, ocrResp)
}

func handleTranslation(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	text, exists := req["text"]
	destLang, exists2 := req["dest_lang"]
	if !exists || !exists2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'text' or 'dest_lang' field"})
		return
	}

	// Gọi service Translation
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{"text": text, "dest_lang": destLang}).
		Post("http://translation_service:5007/translate")
	if err != nil || resp.StatusCode() != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call Translation service"})
		return
	}

	var trResp map[string]string
	json.Unmarshal(resp.Body(), &trResp)

	c.JSON(http.StatusOK, trResp)
}

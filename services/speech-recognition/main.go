package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"context"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ConvertRequest struct {
	AudioURL string `json:"audio_url"`
}

type ConvertResponse struct {
	Text string `json:"text"`
}

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

func main() {
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

	r := gin.Default()
	r.POST("/recognize", handleRecognize)
	r.Run(":5004")
}

func handleRecognize(c *gin.Context) {
	var req ConvertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Insert task vào cơ sở dữ liệu
	var taskID int
	err := dbPool.QueryRow(context.Background(),
		"INSERT INTO tasks (service_name, status, input_data) VALUES ($1, $2, $3) RETURNING id",
		"speech-recognition", "processing", map[string]string{"audio_url": req.AudioURL},
	).Scan(&taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// TODO: Thực hiện nhận diện giọng nói
	// Ví dụ: Tải file audio từ req.AudioURL và sử dụng mô hình ASR để chuyển đổi thành text
	// Ở đây, chúng ta giả lập quá trình chuyển đổi
	time.Sleep(2 * time.Second)
	recognized_text := "Recognized speech text"

	// Cập nhật task
	_, err = dbPool.Exec(context.Background(),
		"UPDATE tasks SET status=$1, output_data=$2, updated_at=NOW() WHERE id=$3",
		"completed", map[string]string{"text": recognized_text}, taskID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database update error"})
		return
	}

	c.JSON(http.StatusOK, ConvertResponse{Text: recognized_text})
}

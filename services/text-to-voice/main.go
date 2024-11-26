package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ConvertRequest struct {
	Text     string  `json:"text" binding:"required"`
	Voice    string  `json:"voice"` // voice type (male/female)
	Language string  `json:"language"`
	Speed    float64 `json:"speed" binding:"min=0.1"`  // speech rate (0.1-3.0)
	Pitch    float64 `json:"pitch"`                    // voice pitch (-10 to 10)
	Volume   float64 `json:"volume" binding:"min=0.1"` // volume (0.1-2.0)
}

type ConvertResponse struct {
	AudioURL string `json:"audio_url"`
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
	log.Println("Starting Text-to-Voice service...")

	// Database connection string
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	log.Printf("Connecting to database at %s\n", dbURL)

	// Connect to database
	var err error
	dbPool, err = pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v\n", err)
	}
	defer dbPool.Close()
	log.Println("Database connected successfully")

	r := gin.Default()

	// CORS config
	log.Println("Configuring CORS...")
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:81", "http://text_to_voice_service:5001", "http://management_api:81", "http://localhost:3000", "https://insight.io.vn"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(config))
	r.Static("/audio", "./audio")
	log.Println("CORS configured successfully")

	// Register routes
	r.POST("/tts", handleConvert)
	r.GET("/voices", handleGetVoices)
	log.Println("Routes registered. Starting server at :5001")
	r.Run(":5001")
}

func handleConvert(c *gin.Context) {
	log.Println("Received /tts request")

	// Parse JSON request
	var req ConvertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid input: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	log.Printf("Request payload: %v\n", req)

	// Set default values
	if req.Language == "" {
		req.Language = "en"
	}
	if req.Speed == 0 {
		req.Speed = 1.0
	}
	if req.Volume == 0 {
		req.Volume = 1.0
	}
	if req.Voice == "" {
		req.Voice = "male" // default voice
	}

	// Validate ranges
	if req.Speed > 3.0 {
		req.Speed = 3.0
	}
	if req.Volume > 2.0 {
		req.Volume = 2.0
	}
	if req.Pitch < -10 {
		req.Pitch = -10
	} else if req.Pitch > 10 {
		req.Pitch = 10
	}

	// Insert task into database
	var taskID int
	log.Println("Inserting task into database...")
	inputData := map[string]interface{}{
		"text":     req.Text,
		"language": req.Language,
		"voice":    req.Voice,
		"speed":    req.Speed,
		"pitch":    req.Pitch,
		"volume":   req.Volume,
	}
	inputJSON, _ := json.Marshal(inputData)

	err := dbPool.QueryRow(context.Background(),
		"INSERT INTO tasks (service_name, status, input_data) VALUES ($1, $2, $3) RETURNING id",
		"text-to-voice", "processing", inputJSON,
	).Scan(&taskID)
	if err != nil {
		log.Printf("Database error during task insertion: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	log.Printf("Task inserted with ID: %d\n", taskID)

	// Create audio directory if not exists
	audioDir := "audio"
	if _, err = os.Stat(audioDir); os.IsNotExist(err) {
		log.Printf("Audio directory %s does not exist. Creating...\n", audioDir)
		err = os.Mkdir(audioDir, 0755)
		if err != nil {
			log.Printf("Failed to create audio directory: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create audio directory"})
			return
		}
		log.Printf("Audio directory %s created successfully\n", audioDir)
	}

	// Convert Text-to-Voice with SSML
	audioPath := fmt.Sprintf("output_%d", taskID)
	ssmlText := fmt.Sprintf(`
		<speak>
			<prosody rate="%g%%" pitch="%g%%" volume="%g%%">
				%s
			</prosody>
		</speak>`,
		req.Speed*100,
		req.Pitch*10+100, // Convert pitch range from -10:10 to 0:200
		req.Volume*100,
		req.Text,
	)

	tts := htgotts.Speech{
		Folder:   audioDir,
		Language: req.Language,
		Handler:  &handlers.MPlayer{},
	}

	log.Printf("Converting text to speech. Output file: %s\n", audioPath)
	filePath, err := tts.CreateSpeechFile(ssmlText, audioPath)
	if err != nil {
		log.Printf("TTS conversion failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "TTS conversion failed"})
		return
	}
	log.Printf("TTS conversion succeeded. File path: %s\n", filePath)

	// Update task status and output_data
	audioURL := fmt.Sprintf("http://localhost:5001/audio/output_%d.mp3", taskID)
	outputData := map[string]string{"audio_url": audioURL}
	outputJSON, _ := json.Marshal(outputData)

	log.Printf("Updating task status to 'completed' with audio URL: %s\n", audioURL)
	_, err = dbPool.Exec(context.Background(),
		"UPDATE tasks SET status=$1, output_data=$2, updated_at=NOW() WHERE id=$3",
		"completed", outputJSON, taskID,
	)
	if err != nil {
		log.Printf("Database update error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database update error"})
		return
	}
	log.Printf("Task ID %d updated successfully\n", taskID)

	// Return result
	log.Printf("Returning response with audio URL: %s\n", audioURL)
	c.JSON(http.StatusOK, ConvertResponse{AudioURL: audioURL})
}

func handleGetVoices(c *gin.Context) {
	// Currently supporting basic voice options
	voices := []map[string]interface{}{
		{
			"id":       "male",
			"name":     "Male Voice",
			"language": "en",
			"gender":   "male",
		},
		{
			"id":       "female",
			"name":     "Female Voice",
			"language": "en",
			"gender":   "female",
		},
	}

	c.JSON(http.StatusOK, voices)
}

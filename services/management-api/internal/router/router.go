package router

import (
	"management-api/internal/config"
	"management-api/internal/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(taskHandler handler.TaskHandler, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://202.92.6.77:3000", "http://localhost:3000", "https://insight.io.vn"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}

	r.Use(cors.New(corsConfig))

	r.Static("/uploads", "/shared/images")
	r.Static("/images", "/shared/images")
	//r.Static("/shared", "/shared/images")

	// Các endpoint tương ứng với từng service
	// TTS endpoints
	r.GET("/tts/languages", taskHandler.GetAvailableLanguages)
	r.GET("/tts/voices/:language", taskHandler.GetAvailableVoices)
	r.POST("/tts", taskHandler.HandleTextToVoice)

	r.POST("/vts", taskHandler.HandleVoiceToText)
	r.POST("/remove-bg", taskHandler.HandleBackgroundRemoval)
	r.POST("/speech-recognition", taskHandler.HandleSpeechRecognition)
	r.POST("/face-recognition", taskHandler.HandleFaceRecognition)
	r.POST("/ocr", taskHandler.HandleOCR)
	r.POST("/translate", taskHandler.HandleTranslation)
	r.POST("/upload-audio", taskHandler.UploadAudio)

	return r
}

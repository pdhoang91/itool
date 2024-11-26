package service

import (
	"management-api/internal/domain"
	"management-api/internal/repository"

	"github.com/go-resty/resty/v2"
)

type TaskService interface {
	GetTaskStatus(id int) (*domain.Task, error)
	GetAllTasks() ([]domain.Task, error)
	//HandleTextToVoice(text, language string) (map[string]string, error)
	//HandleTextToVoice(text, language, voice string, speed, pitch, volume float64) (map[string]string, error)
	HandleTextToVoice(text, language, voice string, speed, pitch, volume float64) (map[string]string, error)
	HandleVoiceToText(audioURL string) (map[string]string, error)
	HandleBackgroundRemoval(imagePath string) (string, error)
	HandleSpeechRecognition(audioURL string) (map[string]string, error)
	HandleFaceRecognition(imagePath string) (map[string]interface{}, error)
	HandleOCR(imagePath string) (map[string]string, error)
	HandleTranslation(text, destLang string) (map[string]string, error)
	UploadAudio(filePath string) (string, error)
}

type taskService struct {
	client *resty.Client
	repo   repository.TaskRepository
}

func NewTaskService(client *resty.Client, repo repository.TaskRepository) TaskService {
	return &taskService{
		client: client,
		repo:   repo,
	}
}

func (s *taskService) GetTaskStatus(id int) (*domain.Task, error) {
	return s.repo.GetTask(id)
}

func (s *taskService) GetAllTasks() ([]domain.Task, error) {
	return s.repo.GetAllTasks()
}

// Các phương thức xử lý các dịch vụ như Text-to-Voice, Voice-to-Text, ...

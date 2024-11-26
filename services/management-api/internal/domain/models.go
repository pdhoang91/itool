package domain

import (
	"encoding/json"
	"time"
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

const (
	TaskStatusProcessing = "processing"
	TaskStatusCompleted  = "completed"
	TaskStatusFailed     = "failed"
)

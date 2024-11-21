package domain

import "time"

type Task struct {
	ID          int       `json:"id"`
	ServiceName string    `json:"service_name"`
	Status      string    `json:"status"`
	InputData   []byte    `json:"input_data"`
	OutputData  []byte    `json:"output_data"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

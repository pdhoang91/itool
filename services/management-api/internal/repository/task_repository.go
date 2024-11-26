// repository/task_repository.go

package repository

import (
	"context"
	"fmt"
	"management-api/internal/config"
	"management-api/internal/domain"

	"github.com/jackc/pgx/v4/pgxpool"
)

type TaskRepository interface {
	GetTask(id int) (*domain.Task, error)
	GetAllTasks() ([]domain.Task, error)
	CreateTask(ctx context.Context, task *domain.Task) error
	UpdateTaskStatus(ctx context.Context, taskID int, status string, outputData interface{}) error
	Close()
}

type taskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(cfg config.DatabaseConfig) (TaskRepository, error) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)
	pool, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		return nil, err
	}
	return &taskRepository{db: pool}, nil
}

func (r *taskRepository) CreateTask(ctx context.Context, task *domain.Task) error {
	query := `
        INSERT INTO tasks (service_name, status, input_data, output_data, created_at, updated_at)
        VALUES ($1, $2, $3, $4, NOW(), NOW())
        RETURNING id`

	return r.db.QueryRow(ctx, query,
		task.ServiceName,
		task.Status,
		task.InputData,
		task.OutputData,
	).Scan(&task.ID)
}

func (r *taskRepository) UpdateTaskStatus(ctx context.Context, taskID int, status string, outputData interface{}) error {
	_, err := r.db.Exec(ctx,
		"UPDATE tasks SET status=$1, output_data=$2, updated_at=NOW() WHERE id=$3",
		status, outputData, taskID,
	)
	return err
}

func (r *taskRepository) GetTask(id int) (*domain.Task, error) {
	var task domain.Task
	err := r.db.QueryRow(context.Background(),
		"SELECT id, service_name, status, input_data, output_data, created_at, updated_at FROM tasks WHERE id=$1",
		id,
	).Scan(&task.ID, &task.ServiceName, &task.Status, &task.InputData, &task.OutputData, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) GetAllTasks() ([]domain.Task, error) {
	rows, err := r.db.Query(context.Background(),
		"SELECT id, service_name, status, input_data, output_data, created_at, updated_at FROM tasks ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var task domain.Task
		err := rows.Scan(&task.ID, &task.ServiceName, &task.Status, &task.InputData, &task.OutputData, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *taskRepository) Close() {
	if r.db != nil {
		r.db.Close()
	}
}

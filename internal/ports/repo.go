package ports

import "github.com/Jashanveer-Singh/todo-go/internal/handlers/domain"

type taskRepo interface {
	CreateTask(task domain.Task) error
	UpdateTask(id string, task domain.Task) error
	DeleteTask(id string) error
	GetTasks() ([]domain.Task, error)
}

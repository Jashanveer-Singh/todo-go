package ports

import (
	"github.com/Jashanveer-Singh/todo-go/internal/errr"
	"github.com/Jashanveer-Singh/todo-go/internal/models"
)

type TaskService interface {
	CreateTask(models.TaskRequestDto) *errr.AppError
	UpdateTask(id string, task models.TaskRequestDto) *errr.AppError
	DeleteTask(id string) *errr.AppError
	GetTasks() ([]models.TaskResponseDto, *errr.AppError)
}

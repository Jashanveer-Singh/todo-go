package ports

import (
	"github.com/Jashanveer-Singh/todo-go/internal/errr"
	"github.com/Jashanveer-Singh/todo-go/internal/models"
)

type TaskRepo interface {
	CreateTask(task models.Task) *errr.AppError
	UpdateTask(id int64, task models.Task) *errr.AppError
	DeleteTask(id int64) *errr.AppError
	GetTasks() ([]models.Task, *errr.AppError)
}

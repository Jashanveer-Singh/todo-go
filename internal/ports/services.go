package ports

import (
	"github.com/Jashanveer-Singh/todo-go/internal/errr"
	"github.com/Jashanveer-Singh/todo-go/internal/models"
)

type TaskService interface {
	CreateTask(taskReq models.TaskRequestDto, userID string) *errr.AppError
	UpdateTask(id string, task models.TaskRequestDto, userID string) *errr.AppError
	DeleteTask(id string, userID string) *errr.AppError
	GetTasks(userID string) ([]models.TaskResponseDto, *errr.AppError)
}

type UserService interface {
	CreateUser(models.UserRequestDto) *errr.AppError
	Login(username, password string) (string, *errr.AppError)
}

package ports

import (
	"github.com/Jashanveer-Singh/todo-go/internal/errr"
	"github.com/Jashanveer-Singh/todo-go/internal/models"
)

type TaskService interface {
	CreateTask(taskReq models.TaskRequestDto, claims models.Claims) *errr.AppError
	UpdateTask(id string, task models.TaskRequestDto, claims models.Claims) *errr.AppError
	DeleteTask(id string, claims models.Claims) *errr.AppError
	GetTasks(claims models.Claims) ([]models.TaskResponseDto, *errr.AppError)
}

type UserService interface {
	CreateUser(models.UserRequestDto) *errr.AppError
}

type AuthService interface {
	Login(username, password string) (token string, appErr *errr.AppError)
}

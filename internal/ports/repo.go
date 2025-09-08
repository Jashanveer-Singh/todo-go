package ports

import (
	"github.com/Jashanveer-Singh/todo-go/internal/errr"
	"github.com/Jashanveer-Singh/todo-go/internal/models"
)

type TaskRepo interface {
	SaveTask(task models.Task) *errr.AppError
	UpdateTask(id int64, task models.Task) *errr.AppError
	DeleteTask(id int64, userID int64) *errr.AppError
	GetTasks(userId int64) ([]models.Task, *errr.AppError)
}

type UserRepo interface {
	GetUserByUsername(username string) (models.User, *errr.AppError)
	CreateUser(user models.User) *errr.AppError
}

// type AuthRepo interface

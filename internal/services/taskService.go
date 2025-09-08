package services

import (
	"net/http"
	"strconv"

	"github.com/Jashanveer-Singh/todo-go/internal/errr"
	"github.com/Jashanveer-Singh/todo-go/internal/models"
	"github.com/Jashanveer-Singh/todo-go/internal/ports"
)

type taskService struct {
	taskRepo ports.TaskRepo
}

func NewTaskService(taskRepo ports.TaskRepo) *taskService {
	return &taskService{
		taskRepo,
	}
}

func (ts *taskService) CreateTask(taskReq models.TaskRequestDto, userId string) *errr.AppError {
	id, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		return &errr.AppError{
			Code:    http.StatusBadRequest,
			Message: "invalid id",
		}
	}
	task := taskReq.ToTask()
	task.Status = 0
	task.UserID = id
	if !task.IsValidTask() {
		return &errr.AppError{
			Message: "Invalid task",
			Code:    http.StatusBadRequest,
		}
	}

	appErr := ts.taskRepo.SaveTask(task)
	if appErr != nil {
		return appErr
	}

	return nil
}

func (ts *taskService) UpdateTask(
	idString string,
	taskReq models.TaskRequestDto,
	userID string,
) *errr.AppError {
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		return &errr.AppError{
			Code:    http.StatusBadRequest,
			Message: "invalid id",
		}
	}

	userIDAsInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return &errr.AppError{
			Code:    http.StatusBadRequest,
			Message: "invalid id",
		}
	}
	if taskReq.Title == "" && taskReq.Desc == "" && !taskReq.IsValidStatus() {
		return &errr.AppError{
			Message: "Invalid task format",
			Code:    http.StatusBadRequest,
		}
	}
	task := taskReq.ToTask()
	task.UserID = userIDAsInt

	appErr := ts.taskRepo.UpdateTask(id, task)
	if appErr != nil {
		return appErr
	}

	return nil
}

func (ts *taskService) DeleteTask(idString string, userID string) *errr.AppError {
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		return &errr.AppError{
			Code:    http.StatusBadRequest,
			Message: "invalid id",
		}
	}

	userIDAsInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return &errr.AppError{
			Code:    http.StatusBadRequest,
			Message: "invalid id",
		}
	}

	appErr := ts.taskRepo.DeleteTask(id, userIDAsInt)
	if appErr != nil {
		return appErr
	}

	return nil
}

func (ts *taskService) GetTasks(userId string) ([]models.TaskResponseDto, *errr.AppError) {
	id, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		return []models.TaskResponseDto{}, &errr.AppError{
			Code:    http.StatusBadRequest,
			Message: "invalid id",
		}
	}
	tasks, appErr := ts.taskRepo.GetTasks(id)
	if appErr != nil {
		return nil, appErr
	}

	taskRes := make([]models.TaskResponseDto, len(tasks))

	for i := range tasks {
		taskRes[i] = tasks[i].ToDto()
	}

	return taskRes, nil
}

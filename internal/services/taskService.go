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
	// authService ports.AuthService
}

func NewTaskService(taskRepo ports.TaskRepo) *taskService {
	return &taskService{
		taskRepo: taskRepo,
		// authService: authService,
	}
}

func (ts *taskService) CreateTask(
	taskReq models.TaskRequestDto,
	claims models.Claims,
) *errr.AppError {
	task := taskReq.ToTask()
	task.Status = 2
	task.UserID = claims.ID
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
	taskIDStr string,
	taskReq models.TaskRequestDto,
	claims models.Claims,
) *errr.AppError {
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		return errr.NewBadRequestError("Invalid task id")
	}

	if taskReq.Title == "" && taskReq.Desc == "" && !taskReq.IsValidStatus() {
		return &errr.AppError{
			Message: "Invalid task format",
			Code:    http.StatusBadRequest,
		}
	}
	task := taskReq.ToTask()
	task.UserID = claims.ID

	appErr := ts.taskRepo.UpdateTask(taskID, task)
	if appErr != nil {
		return appErr
	}

	return nil
}

func (ts *taskService) DeleteTask(idString string, claims models.Claims) *errr.AppError {
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		return &errr.AppError{
			Code:    http.StatusBadRequest,
			Message: "invalid id",
		}
	}

	appErr := ts.taskRepo.DeleteTask(id, claims.ID)
	if appErr != nil {
		return appErr
	}

	return nil
}

func (ts *taskService) GetTasks(claims models.Claims) ([]models.TaskResponseDto, *errr.AppError) {
	tasks, appErr := ts.taskRepo.GetTasks(claims.ID)
	if appErr != nil {
		return nil, appErr
	}

	taskRes := make([]models.TaskResponseDto, len(tasks))

	for i := range tasks {
		taskRes[i] = tasks[i].ToDto()
	}

	return taskRes, nil
}

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

func (ts *taskService) CreateTask(taskReq models.TaskRequestDto) *errr.AppError {
	task := taskReq.ToTask()
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

func (ts *taskService) UpdateTask(idString string, taskReq models.TaskRequestDto) *errr.AppError {
	id, err := strconv.ParseInt(idString, 10, 64)
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

	appErr := ts.taskRepo.UpdateTask(id, task)
	if appErr != nil {
		return appErr
	}

	return nil
}

func (ts *taskService) DeleteTask(idString string) *errr.AppError {
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		return &errr.AppError{
			Code:    http.StatusBadRequest,
			Message: "invalid id",
		}
	}

	appErr := ts.taskRepo.DeleteTask(id)
	if appErr != nil {
		return appErr
	}

	return nil
}

func (ts *taskService) GetTasks() ([]models.TaskResponseDto, *errr.AppError) {
	tasks, appErr := ts.taskRepo.GetTasks()
	if appErr != nil {
		return nil, appErr
	}

	taskRes := make([]models.TaskResponseDto, len(tasks))

	for i := range tasks {
		taskRes[i] = tasks[i].ToDto()
	}

	return taskRes, nil
}

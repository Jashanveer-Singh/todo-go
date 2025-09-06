package ports

import "github.com/Jashanveer-Singh/todo-go/internal/models"

type SendResponse = func(Response)

type apiServer interface {
	ListenCreateTaskRequest() (models.TaskRequestDto, SendResponse)
	ListenUpdatetaskRequest() (models.TaskRequestDto, SendResponse)
	ListenDeleteTaskRequest() (models.TaskRequestDto, SendResponse)
	ListenGettaskRequest() SendResponse
}

type Response struct {
	StatusCode int
	Body       string
}

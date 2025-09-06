package ports

import "net/http"

type TaskHandlers interface {
	CreateTaskHandler(w http.ResponseWriter, r *http.Request)
	UpdateTaskHandler(w http.ResponseWriter, r *http.Request)
	DeleteTaskHandler(w http.ResponseWriter, r *http.Request)
	GetTasksHandler(w http.ResponseWriter, r *http.Request)
}

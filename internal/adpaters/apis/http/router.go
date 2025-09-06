package http

import (
	"net/http"

	"github.com/Jashanveer-Singh/todo-go/internal/ports"
)

func newRouter(handler ports.TaskHandlers) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tasks", handler.GetTasksHandler)
	mux.HandleFunc("POST /tasks", handler.CreateTaskHandler)
	mux.HandleFunc("PUT /tasks/{id}", handler.UpdateTaskHandler)
	mux.HandleFunc("DELETE /tasks/{id}", handler.DeleteTaskHandler)

	return mux
}

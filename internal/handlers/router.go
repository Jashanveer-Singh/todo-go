package handlers

import (
	"net/http"

	"github.com/Jashanveer-Singh/todo-go/internal/ports"
)

func NewRouter(handler ports.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tasks", handler.GetTasks)
	mux.HandleFunc("POST /tasks", handler.CreateTask)
	mux.HandleFunc("PUT /tasks/{id}", handler.UpdateTask)
	mux.HandleFunc("DELETE /tasks/{id}", handler.DeleteTask)

	return mux
}

package http

import (
	"net/http"
)

func newRouter(taskHandler *taskHandler, userHandler *userHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tasks", taskHandler.GetTasksHandler)
	mux.HandleFunc("POST /tasks", taskHandler.CreateTaskHandler)
	mux.HandleFunc("PUT /tasks/{id}", taskHandler.UpdateTaskHandler)
	mux.HandleFunc("DELETE /tasks/{id}", taskHandler.DeleteTaskHandler)

	mux.HandleFunc("POST /users", userHandler.CreateUserHandler)
	mux.HandleFunc("POST /users/login", userHandler.Login)

	return mux
}

package http

import (
	"net/http"
)

func newRouter(
	taskHandler *taskHandler,
	userHandler *userHandler,
	authHandler *authHandler,
	authMiddleware *AuthMiddleware,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc(
		"GET /tasks",
		authMiddleware.isAuthenticatedMiddleware(taskHandler.GetTasksHandler),
	)
	mux.HandleFunc(
		"POST /tasks",
		authMiddleware.isAuthenticatedMiddleware(taskHandler.CreateTaskHandler),
	)
	mux.HandleFunc(
		"PUT /tasks/{id}",
		authMiddleware.isAuthenticatedMiddleware(taskHandler.UpdateTaskHandler),
	)
	mux.HandleFunc(
		"DELETE /tasks/{id}",
		authMiddleware.isAuthenticatedMiddleware(taskHandler.DeleteTaskHandler),
	)

	mux.HandleFunc("POST /users", userHandler.CreateUserHandler)
	mux.HandleFunc("POST /auth", authHandler.Login)

	return mux
}

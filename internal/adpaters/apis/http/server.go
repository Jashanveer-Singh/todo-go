package http

import (
	"net/http"

	"github.com/Jashanveer-Singh/todo-go/internal/ports"
)

func NewHttpServer(
	taskService ports.TaskService,
	userService ports.UserService,
	authService ports.AuthService,
	tokenProvider ports.TokenProvider,
) httpServer {
	return httpServer{
		taskService:   taskService,
		userService:   userService,
		tokenProvider: tokenProvider,
		authService:   authService,
	}
}

type httpServer struct {
	taskService   ports.TaskService
	userService   ports.UserService
	tokenProvider ports.TokenProvider
	authService   ports.AuthService
}

func (hs httpServer) ListenAndServe(addr string) {
	taskHandler := newTaskHandler(hs.taskService)
	userHandler := NewUserHandler(hs.userService)
	authHandler := NewAuthHandler(hs.authService)
	authMiddleware := NewAuthMiddleware(hs.tokenProvider)
	router := newRouter(taskHandler, userHandler, authHandler, authMiddleware)
	http.ListenAndServe(addr, router)
}

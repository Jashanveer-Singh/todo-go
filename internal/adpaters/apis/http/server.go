package http

import (
	"net/http"

	"github.com/Jashanveer-Singh/todo-go/internal/ports"
)

func NewHttpServer(ts ports.TaskService, us ports.UserService) httpServer {
	return httpServer{
		ts: ts,
		us: us,
	}
}

type httpServer struct {
	ts ports.TaskService
	us ports.UserService
}

func (hs httpServer) ListenAndServe(addr string) {
	taskHandler := newTaskHandler(hs.ts)
	userHandler := NewUserHandler(hs.us)
	router := newRouter(taskHandler, userHandler)
	http.ListenAndServe(addr, router)
}

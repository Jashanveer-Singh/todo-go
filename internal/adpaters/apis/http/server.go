package http

import (
	"net/http"

	"github.com/Jashanveer-Singh/todo-go/internal/services"
)

func NewHttpServer(ts services.TaskService) httpServer {
	return httpServer{
		ts,
	}
}

type httpServer struct {
	ts services.TaskService
}

func (hs httpServer) ListenAndServe(addr string) {
	handlers := newTaskHandler(hs.ts)
	router := newRouter(handlers)
	http.ListenAndServe(addr, router)
}

package http

import (
	"net/http"

	"github.com/Jashanveer-Singh/todo-go/internal/ports"
)

func NewHttpServer(ts ports.TaskService) httpServer {
	return httpServer{
		ts,
	}
}

type httpServer struct {
	ts ports.TaskService
}

func (hs httpServer) ListenAndServe(addr string) {
	handlers := newTaskHandler(hs.ts)
	router := newRouter(handlers)
	http.ListenAndServe(addr, router)
}

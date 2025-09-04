package apis

import "net/http"

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tasks")
	mux.HandleFunc("POST /tasks")
	mux.HandleFunc("PUT /tasks")
	mux.HandleFunc("DELETE /tasks")

	return mux
}

package http

import (
	"encoding/json"
	"net/http"

	"github.com/Jashanveer-Singh/todo-go/internal/models"
	"github.com/Jashanveer-Singh/todo-go/internal/ports"
)

type taskHandler struct {
	ts ports.TaskService
}

func newTaskHandler(ts ports.TaskService) *taskHandler {
	return &taskHandler{
		ts,
	}
}

func (th taskHandler) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(models.Claims)
	if !ok {
		http.Error(w, "Unexpected error in authentication", http.StatusInternalServerError)
	}
	taskRes, appErr := th.ts.GetTasks(claims)

	if appErr != nil {
		http.Error(w, appErr.Message, appErr.Code)
		return
	}

	tasksjson, _ := json.Marshal(taskRes)

	w.Header().Set("content", "application/json")
	w.Write(tasksjson)
}

func (th taskHandler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(models.Claims)
	if !ok {
		http.Error(w, "Unexpected error in authentication", http.StatusInternalServerError)
	}
	id := r.PathValue("id")

	var taskReq models.TaskRequestDto

	err := json.NewDecoder(r.Body).Decode(&taskReq)
	if err != nil {
		http.Error(w, "Invalid Body", http.StatusBadRequest)
		return
	}

	appErr := th.ts.UpdateTask(id, taskReq, claims)

	if appErr != nil {
		http.Error(w, appErr.Message, appErr.Code)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(""))
}

func (th taskHandler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(models.Claims)
	if !ok {
		http.Error(w, "Unexpected error in authentication", http.StatusInternalServerError)
	}
	var taskReq models.TaskRequestDto
	err := json.NewDecoder(r.Body).Decode(&taskReq)
	if err != nil {
		http.Error(w, "Invalid Body", http.StatusBadRequest)
		return
	}

	appErr := th.ts.CreateTask(taskReq, claims)
	if appErr != nil {
		http.Error(w, appErr.Message, appErr.Code)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("task created successfully"))
}

func (th taskHandler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(models.Claims)
	if !ok {
		http.Error(w, "Unexpected error in authentication", http.StatusInternalServerError)
	}
	id := r.PathValue("id")
	appErr := th.ts.DeleteTask(id, claims)
	if appErr != nil {
		http.Error(w, appErr.Message, appErr.Code)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(""))
}

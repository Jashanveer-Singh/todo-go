package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Jashanveer-Singh/todo-go/internal/models"
	"github.com/Jashanveer-Singh/todo-go/internal/ports"
)

type handler struct {
	repo ports.TaskRepo
}

func NewHandler(repo ports.TaskRepo) *handler {
	return &handler{
		repo: repo,
	}
}

func (h handler) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, appErr := h.repo.GetTasks()
	if appErr != nil {
		http.Error(w, appErr.Message, 500)
		return
	}

	taskRes := make([]models.TaskResponseDto, len(tasks))

	for i := range tasks {
		taskRes[i] = tasks[i].ToDto()
	}

	tasksjson, err := json.Marshal(taskRes)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content", "application/json")
	w.Write(tasksjson)
}

func (h handler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	var taskReq models.TaskRequestDto

	err = json.NewDecoder(r.Body).Decode(&taskReq)
	if err != nil {
		http.Error(w, "Invalid Body", http.StatusBadRequest)
	}

	task := taskReq.ToTask()

	appErr := h.repo.UpdateTask(id, task)
	if err != nil {
		http.Error(w, appErr.Message, http.StatusInternalServerError)
	}

	w.Write([]byte(""))
}

func (h handler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var taskRequest models.TaskRequestDto
	err := json.NewDecoder(r.Body).Decode(&taskRequest)
	if err != nil {
		http.Error(w, "Invalid request Body", http.StatusBadRequest)
		return
	}

	task := taskRequest.ToTask()

	appErr := h.repo.SaveTask(task)
	if appErr != nil {
		http.Error(w, appErr.Message, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(""))
}

func (h handler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid Id", http.StatusBadRequest)
		return
	}

	appErr := h.repo.DeleteTask(id)
	if appErr != nil {
		http.Error(w, appErr.Message, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(""))
}

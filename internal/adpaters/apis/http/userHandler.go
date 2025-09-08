package http

import (
	"encoding/json"
	"net/http"

	"github.com/Jashanveer-Singh/todo-go/internal/models"
	"github.com/Jashanveer-Singh/todo-go/internal/ports"
)

type userHandler struct {
	us ports.UserService
}

func NewUserHandler(us ports.UserService) *userHandler {
	return &userHandler{
		us: us,
	}
}

func (uh userHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	userReq := models.UserRequestDto{}

	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		http.Error(w, "Invalid Body", http.StatusBadRequest)
		return
	}

	appErr := uh.us.CreateUser(userReq)
	if appErr != nil {
		http.Error(w, appErr.Message, appErr.Code)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User Created Successfully"))
}

func (uh userHandler) Login(w http.ResponseWriter, r *http.Request) {
	userReq := models.UserRequestDto{}

	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		http.Error(w, "Invalid Body", http.StatusBadRequest)
		return
	}

	id, appErr := uh.us.Login(userReq.Username, userReq.Password)
	if appErr != nil {
		http.Error(w, appErr.Message, appErr.Code)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(id))
}

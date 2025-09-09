package http

import (
	"encoding/json"
	"net/http"

	"github.com/Jashanveer-Singh/todo-go/internal/models"
	"github.com/Jashanveer-Singh/todo-go/internal/ports"
)

func NewAuthHandler(authService ports.AuthService) *authHandler {
	return &authHandler{
		authService: authService,
	}
}

type authHandler struct {
	authService ports.AuthService
}

func (ah authHandler) Login(w http.ResponseWriter, r *http.Request) {
	userReq := models.UserRequestDto{}

	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		http.Error(w, "Invalid Body", http.StatusBadRequest)
		return
	}

	token, appErr := ah.authService.Login(userReq.Username, userReq.Password)
	if appErr != nil {
		http.Error(w, appErr.Message, appErr.Code)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}

package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Jashanveer-Singh/todo-go/internal/errr"
	"github.com/Jashanveer-Singh/todo-go/internal/models"
	"github.com/Jashanveer-Singh/todo-go/test/mocks"
	"github.com/golang/mock/gomock"
)

func Test_userHandler_CreateUserHandler(t *testing.T) {
	tests := []struct {
		name         string // description of this test case
		setupMUS     func(mus *mocks.MockUserService)
		requestBody  io.Reader
		wantStatus   int
		responseBody string
	}{
		{
			name:         "invalid request body",
			setupMUS:     func(mus *mocks.MockUserService) {},
			requestBody:  strings.NewReader("adsfj;lsdj"),
			wantStatus:   http.StatusBadRequest,
			responseBody: "Invalid Body\n",
		},
		{
			name: "user service returns error",
			setupMUS: func(mus *mocks.MockUserService) {
				mus.EXPECT().CreateUser(models.UserRequestDto{
					Username: "jass",
					Password: "password",
				}).Return(&errr.AppError{
					Code:    http.StatusInternalServerError,
					Message: "error message from user service",
				})
			},
			requestBody:  strings.NewReader(`{"username": "jass", "password": "password"}`),
			wantStatus:   http.StatusInternalServerError,
			responseBody: "error message from user service\n",
		},
		{
			name: "successful response",
			setupMUS: func(mus *mocks.MockUserService) {
				mus.EXPECT().CreateUser(models.UserRequestDto{
					Username: "jass",
					Password: "password",
				}).Return(nil)
			},
			requestBody:  strings.NewReader(`{"username": "jass", "password": "password"}`),
			wantStatus:   http.StatusCreated,
			responseBody: "User Created Successfully",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/tasks/asdf", tt.requestBody)
			rr := httptest.NewRecorder()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserService := mocks.NewMockUserService(ctrl)
			tt.setupMUS(mockUserService)
			uh := NewUserHandler(mockUserService)
			uh.CreateUserHandler(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("wanted status code %d, got %d.", tt.wantStatus, rr.Code)
			}
			if rr.Body.String() != tt.responseBody {
				t.Errorf("wanted response body: %s, got %s.", tt.responseBody, rr.Body)
			}
		})
	}
}

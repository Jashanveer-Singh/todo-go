package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Jashanveer-Singh/todo-go/internal/errr"
	"github.com/Jashanveer-Singh/todo-go/test/mocks"
	"github.com/golang/mock/gomock"
)

func Test_authHandler_Login(t *testing.T) {
	tests := []struct {
		name         string // description of this test case
		setupMAS     func(mus *mocks.MockAuthService)
		requestBody  io.Reader
		wantStatus   int
		responseBody string
	}{
		{
			name:         "invalid request body",
			setupMAS:     func(mus *mocks.MockAuthService) {},
			requestBody:  strings.NewReader("adsfj;lsdj"),
			wantStatus:   http.StatusBadRequest,
			responseBody: "Invalid Body\n",
		},
		{
			name: "user service returns error",
			setupMAS: func(mus *mocks.MockAuthService) {
				mus.EXPECT().Login("jass", "password").
					Return("", &errr.AppError{
						Code:    http.StatusInternalServerError,
						Message: "error message from auth service",
					})
			},
			requestBody:  strings.NewReader(`{"username": "jass", "password": "password"}`),
			wantStatus:   http.StatusInternalServerError,
			responseBody: "error message from auth service\n",
		},
		{
			name: "successful response",
			setupMAS: func(mus *mocks.MockAuthService) {
				mus.EXPECT().Login("jass", "password").Return("token", nil)
			},
			requestBody:  strings.NewReader(`{"username": "jass", "password": "password"}`),
			wantStatus:   http.StatusOK,
			responseBody: "token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/tasks/asdf", tt.requestBody)
			rr := httptest.NewRecorder()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthService := mocks.NewMockAuthService(ctrl)
			tt.setupMAS(mockAuthService)
			uh := NewAuthHandler(mockAuthService)
			uh.Login(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("wanted status code %d, got %d.", tt.wantStatus, rr.Code)
			}
			if rr.Body.String() != tt.responseBody {
				t.Errorf("wanted response body: %s, got %s.", tt.responseBody, rr.Body)
			}
		})
	}
}

package http

import (
	"context"
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

func Test_taskHandler_GetTasksHandler(t *testing.T) {
	tests := []struct {
		name     string
		setupMTS func(*mocks.MockTaskService)
		// url          string
		// requestBody  io.Reader
		wantStatus   int
		responseBody string
	}{
		{
			name: "successful response",
			setupMTS: func(mts *mocks.MockTaskService) {
				mts.EXPECT().GetTasks(models.Claims{ID: 4321}).Return([]models.TaskResponseDto{
					{
						ID:     "1234",
						Title:  "title",
						Desc:   "desc",
						Status: "Pending",
					},
				}, nil)
			},
			wantStatus:   http.StatusOK,
			responseBody: `[{"id":"1234","title":"title","desc":"desc","status":"Pending"}]`,
		},
		{
			name: "task service get task returns error",
			setupMTS: func(mts *mocks.MockTaskService) {
				mts.EXPECT().GetTasks(models.Claims{ID: 4321}).Return(nil, &errr.AppError{
					Code:    http.StatusInternalServerError,
					Message: "error message",
				})
			},
			wantStatus:   http.StatusInternalServerError,
			responseBody: "error message\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
			req = req.WithContext(context.WithValue(req.Context(), "claims", models.Claims{
				ID: 4321,
			}))
			rr := httptest.NewRecorder()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTaskService := mocks.NewMockTaskService(ctrl)
			tt.setupMTS(mockTaskService)
			th := newTaskHandler(mockTaskService)
			th.GetTasksHandler(rr, req)
			if rr.Code != tt.wantStatus {
				t.Errorf("wanted status code %d, got %d.", tt.wantStatus, rr.Code)
			}
			if rr.Body.String() != tt.responseBody {
				t.Errorf("wanted response body: %s, got %s.", tt.responseBody, rr.Body)
			}
		})
	}
}

func Test_taskHandler_UpdateTaskHandler(t *testing.T) {
	tests := []struct {
		name         string
		setupMTS     func(*mocks.MockTaskService)
		url          string
		requestBody  io.Reader
		wantStatus   int
		responseBody string
	}{
		{
			name: "successful response",
			setupMTS: func(mts *mocks.MockTaskService) {
				mts.EXPECT().UpdateTask("1234324", models.TaskRequestDto{
					Title:  "title",
					Desc:   "desc",
					Status: "Pending",
				}, models.Claims{ID: 4321}).Return(nil)
			},
			url: "/tasks/1234324",
			requestBody: strings.NewReader(`{
				"title": "title",
				"desc": "desc",
				"status": "Pending"
				}`),
			wantStatus:   http.StatusNoContent,
			responseBody: "",
		},
		{
			name: "only title is updated",
			setupMTS: func(mts *mocks.MockTaskService) {
				mts.EXPECT().UpdateTask("1234324", models.TaskRequestDto{
					Title: "title",
				}, models.Claims{ID: 4321}).Return(nil)
			},
			url:          "/tasks/1234324",
			requestBody:  strings.NewReader(`{"title": "title"}`),
			wantStatus:   http.StatusNoContent,
			responseBody: "",
		},
		{
			name: "only description is updated",
			setupMTS: func(mts *mocks.MockTaskService) {
				mts.EXPECT().UpdateTask("1234324", models.TaskRequestDto{
					Desc: "desc",
				}, models.Claims{ID: 4321}).Return(nil)
			},
			url:          "/tasks/1234324",
			requestBody:  strings.NewReader(`{"desc": "desc"}`),
			wantStatus:   http.StatusNoContent,
			responseBody: "",
		},
		{
			name:         "invalid task update request",
			setupMTS:     func(mts *mocks.MockTaskService) {},
			url:          "/tasks/1234324",
			requestBody:  strings.NewReader(`{"desc": 123}`),
			wantStatus:   http.StatusBadRequest,
			responseBody: "Invalid Body\n",
		},
		{
			name: "task service returns error",
			setupMTS: func(mts *mocks.MockTaskService) {
				mts.EXPECT().
					UpdateTask("123", models.TaskRequestDto{}, models.Claims{ID: 4321}).
					Return(&errr.AppError{
						Code:    http.StatusBadGateway,
						Message: "error message",
					})
			},
			url:          "/tasks/123",
			requestBody:  strings.NewReader("{}"),
			wantStatus:   http.StatusBadGateway,
			responseBody: "error message\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPut, tt.url, tt.requestBody)
			req.Header.Set("Authorization", "Bearer token")
			rr := httptest.NewRecorder()
			taskServiceCtrl := gomock.NewController(t)
			defer taskServiceCtrl.Finish()
			mockTaskService := mocks.NewMockTaskService(taskServiceCtrl)
			tt.setupMTS(mockTaskService)

			tokenProviderCtrl := gomock.NewController(t)
			defer tokenProviderCtrl.Finish()
			mockTokenProvider := mocks.NewMockTokenProvider(tokenProviderCtrl)
			mockTokenProvider.EXPECT().ValidateToken("token").Return(models.Claims{ID: 4321}, nil)

			am := NewAuthMiddleware(mockTokenProvider)
			th := newTaskHandler(mockTaskService)
			uh := NewUserHandler(nil)
			ah := NewAuthHandler(nil)
			router := newRouter(th, uh, ah, am)
			router.ServeHTTP(rr, req)
			if rr.Code != tt.wantStatus {
				t.Errorf("wanted status code %d, got %d.", tt.wantStatus, rr.Code)
			}
			if rr.Body.String() != tt.responseBody {
				t.Errorf("wanted response body: %s, got %s.", tt.responseBody, rr.Body)
			}
		})
	}
}

func Test_taskHandler_DeleteTaskHandler(t *testing.T) {
	tests := []struct {
		name         string
		setupMTS     func(*mocks.MockTaskService)
		url          string
		wantStatus   int
		responseBody string
	}{
		{
			name: "successful response",
			setupMTS: func(mts *mocks.MockTaskService) {
				mts.EXPECT().DeleteTask("1234", models.Claims{ID: 4321}).Return(nil)
			},
			url:          "/tasks/1234",
			wantStatus:   http.StatusNoContent,
			responseBody: "",
		},
		{
			name: "task service returns error",
			setupMTS: func(mts *mocks.MockTaskService) {
				mts.EXPECT().DeleteTask("1234", models.Claims{ID: 4321}).Return(&errr.AppError{
					Code:    http.StatusBadGateway,
					Message: "error message",
				})
			},
			url:          "/tasks/1234",
			wantStatus:   http.StatusBadGateway,
			responseBody: "error message\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, tt.url, nil)
			req.Header.Set("Authorization", "Bearer token")
			rr := httptest.NewRecorder()
			taskServiceCtrl := gomock.NewController(t)
			defer taskServiceCtrl.Finish()
			mockTaskService := mocks.NewMockTaskService(taskServiceCtrl)
			tt.setupMTS(mockTaskService)

			tokenProviderCtrl := gomock.NewController(t)
			defer tokenProviderCtrl.Finish()
			mockTokenProvider := mocks.NewMockTokenProvider(tokenProviderCtrl)
			mockTokenProvider.EXPECT().ValidateToken("token").Return(models.Claims{ID: 4321}, nil)

			am := NewAuthMiddleware(mockTokenProvider)
			th := newTaskHandler(mockTaskService)
			uh := NewUserHandler(nil)
			ah := NewAuthHandler(nil)
			router := newRouter(th, uh, ah, am)
			router.ServeHTTP(rr, req)
			if rr.Code != tt.wantStatus {
				t.Errorf("wanted status code %d, got %d.", tt.wantStatus, rr.Code)
			}
			if rr.Body.String() != tt.responseBody {
				t.Errorf("wanted response body: %s, got %s.", tt.responseBody, rr.Body)
			}
		})
	}
}

func Test_taskHandler_CreateTaskHandler(t *testing.T) {
	tests := []struct {
		name         string
		setupMTS     func(*mocks.MockTaskService)
		requestBody  io.Reader
		wantStatus   int
		responseBody string
	}{
		{
			name: "successful response",
			setupMTS: func(mts *mocks.MockTaskService) {
				mts.EXPECT().CreateTask(models.TaskRequestDto{
					Title: "title",
					Desc:  "desc",
				}, models.Claims{ID: 4321}).Return(nil)
			},
			requestBody:  strings.NewReader(`{"title":"title","desc":"desc"}`),
			wantStatus:   http.StatusCreated,
			responseBody: "task created successfully",
		},
		{
			name: "task service returns error",
			setupMTS: func(mts *mocks.MockTaskService) {
				mts.EXPECT().
					CreateTask(models.TaskRequestDto{}, models.Claims{ID: 4321}).
					Return(&errr.AppError{
						Code:    http.StatusBadGateway,
						Message: "error message",
					})
			},
			requestBody:  strings.NewReader("{}"),
			wantStatus:   http.StatusBadGateway,
			responseBody: "error message\n",
		},
		{
			name:         "invalid request body",
			setupMTS:     func(mts *mocks.MockTaskService) {},
			requestBody:  strings.NewReader(`{"lod":'`),
			wantStatus:   http.StatusBadRequest,
			responseBody: "Invalid Body\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/tasks/asdf", tt.requestBody)
			req = req.WithContext(context.WithValue(req.Context(), "claims", models.Claims{
				ID: 4321,
			}))
			rr := httptest.NewRecorder()
			taskServiceCtrl := gomock.NewController(t)
			defer taskServiceCtrl.Finish()

			mockTaskService := mocks.NewMockTaskService(taskServiceCtrl)
			tt.setupMTS(mockTaskService)
			th := newTaskHandler(mockTaskService)
			th.CreateTaskHandler(rr, req)
			if rr.Code != tt.wantStatus {
				t.Errorf("wanted status code %d, got %d.", tt.wantStatus, rr.Code)
			}
			if rr.Body.String() != tt.responseBody {
				t.Errorf("wanted response body: %s, got %s.", tt.responseBody, rr.Body)
			}
		})
	}
}

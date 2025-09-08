package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Jashanveer-Singh/todo-go/internal/errr"
	"github.com/Jashanveer-Singh/todo-go/internal/models"
	"github.com/Jashanveer-Singh/todo-go/internal/ports/mocks"
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
				mts.EXPECT().GetTasks().Return([]models.TaskResponseDto{
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
				mts.EXPECT().GetTasks().Return(nil, &errr.AppError{
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
				mts.EXPECT().UpdateTask("", models.TaskRequestDto{
					Title:  "title",
					Desc:   "desc",
					Status: "Pending",
				}).Return(nil)
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
				mts.EXPECT().UpdateTask("", models.TaskRequestDto{
					Title: "title",
				}).Return(nil)
			},
			url:          "/tasks/1234324",
			requestBody:  strings.NewReader(`{"title": "title"}`),
			wantStatus:   http.StatusNoContent,
			responseBody: "",
		},
		{
			name: "only description is updated",
			setupMTS: func(mts *mocks.MockTaskService) {
				mts.EXPECT().UpdateTask("", models.TaskRequestDto{
					Desc: "desc",
				}).Return(nil)
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
				mts.EXPECT().UpdateTask("", models.TaskRequestDto{}).Return(&errr.AppError{
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
			req := httptest.NewRequest(http.MethodGet, tt.url, tt.requestBody)
			rr := httptest.NewRecorder()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTaskService := mocks.NewMockTaskService(ctrl)
			tt.setupMTS(mockTaskService)
			th := newTaskHandler(mockTaskService)
			th.UpdateTaskHandler(rr, req)
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
		wantStatus   int
		responseBody string
	}{
		{
			name: "successful response",
			setupMTS: func(mts *mocks.MockTaskService) {
				mts.EXPECT().DeleteTask("").Return(nil)
			},
			wantStatus:   http.StatusNoContent,
			responseBody: "",
		},
		{
			name: "task service returns error",
			setupMTS: func(mts *mocks.MockTaskService) {
				mts.EXPECT().DeleteTask("").Return(&errr.AppError{
					Code:    http.StatusBadGateway,
					Message: "error message",
				})
			},
			wantStatus:   http.StatusBadGateway,
			responseBody: "error message\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/tasks/asdf", nil)
			rr := httptest.NewRecorder()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTaskService := mocks.NewMockTaskService(ctrl)
			tt.setupMTS(mockTaskService)
			th := newTaskHandler(mockTaskService)
			th.DeleteTaskHandler(rr, req)
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
				}).Return(nil)
			},
			requestBody:  strings.NewReader(`{"title":"title","desc":"desc"}`),
			wantStatus:   http.StatusCreated,
			responseBody: "",
		},
		{
			name: "task service returns error",
			setupMTS: func(mts *mocks.MockTaskService) {
				mts.EXPECT().CreateTask(models.TaskRequestDto{}).Return(&errr.AppError{
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
			rr := httptest.NewRecorder()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTaskService := mocks.NewMockTaskService(ctrl)
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

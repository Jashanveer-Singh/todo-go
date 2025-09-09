package services

import (
	"net/http"
	"slices"
	"testing"

	"github.com/Jashanveer-Singh/todo-go/internal/errr"
	"github.com/Jashanveer-Singh/todo-go/internal/models"
	"github.com/Jashanveer-Singh/todo-go/test/mocks"
	"github.com/golang/mock/gomock"
)

func Test_taskService_CreateTask(t *testing.T) {
	tests := []struct {
		name          string
		setupTaskRepo func(mtr *mocks.MockTaskRepo)
		taskReq       models.TaskRequestDto
		appErr        *errr.AppError
		claims        models.Claims
	}{
		{
			name: "successfully created task",
			setupTaskRepo: func(mtr *mocks.MockTaskRepo) {
				mtr.EXPECT().SaveTask(models.Task{
					Title:  "title",
					Desc:   "desc",
					Status: 2,
					UserID: 1234,
				}).Return(nil)
			},
			taskReq: models.TaskRequestDto{
				Title:  "title",
				Desc:   "desc",
				Status: "Pending",
			},
			appErr: nil,
			claims: models.Claims{
				ID:   1234,
				Role: "",
			},
		},
		{
			name: "successfully created task with invalid status input",
			setupTaskRepo: func(mtr *mocks.MockTaskRepo) {
				mtr.EXPECT().SaveTask(models.Task{
					Title:  "title",
					Desc:   "desc",
					Status: 2,
					UserID: 1234,
				}).Return(nil)
			},
			taskReq: models.TaskRequestDto{
				Title:  "title",
				Desc:   "desc",
				Status: "sdaf",
			},
			appErr: nil,
			claims: models.Claims{
				ID:   1234,
				Role: "",
			},
		},
		{
			name:          "failed to create task because of empty title",
			setupTaskRepo: func(mtr *mocks.MockTaskRepo) {},
			taskReq: models.TaskRequestDto{
				Title:  "",
				Desc:   "desc",
				Status: "Pending",
			},
			appErr: &errr.AppError{
				Message: "Invalid task",
				Code:    http.StatusBadRequest,
			},
			claims: models.Claims{
				ID:   1234,
				Role: "",
			},
		},
		{
			name:          "failed to create task because of empty desc",
			setupTaskRepo: func(mtr *mocks.MockTaskRepo) {},
			taskReq: models.TaskRequestDto{
				Title:  "title",
				Desc:   "",
				Status: "Pending",
			},
			appErr: &errr.AppError{
				Message: "Invalid task",
				Code:    http.StatusBadRequest,
			},
			claims: models.Claims{
				ID:   1234,
				Role: "",
			},
		},
		{
			name: "task repo failed to save task",
			setupTaskRepo: func(mtr *mocks.MockTaskRepo) {
				mtr.EXPECT().SaveTask(models.Task{
					Title:  "title",
					Desc:   "desc",
					Status: 2,
					UserID: 1234,
				}).Return(&errr.AppError{
					Code:    0,
					Message: "error message from task repo",
				})
			},
			taskReq: models.TaskRequestDto{
				Title:  "title",
				Desc:   "desc",
				Status: "Pending",
			},
			appErr: &errr.AppError{
				Code:    0,
				Message: "error message from task repo",
			},
			claims: models.Claims{
				ID:   1234,
				Role: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mtr := mocks.NewMockTaskRepo(ctrl)
			tt.setupTaskRepo(mtr)
			ts := NewTaskService(mtr)

			got := ts.CreateTask(tt.taskReq, tt.claims)
			if tt.appErr == nil && tt.appErr != got {
				t.Errorf("CreateTask() failed, got err: %v.", got)
				return
			}
			if tt.appErr != nil && got == nil {
				t.Errorf("CreateTask() successed unexpectedly, wanted err: %v.", tt.appErr)
				return
			}
			if tt.appErr != nil && *tt.appErr != *got {
				t.Errorf("CreateTask() = %v, want %v", got, tt.appErr)
			}
		})
	}
}

func Test_taskService_UpdateTask(t *testing.T) {
	tests := []struct {
		name          string
		setupTaskRepo func(mtr *mocks.MockTaskRepo)
		id            string
		taskReq       models.TaskRequestDto
		appErr        *errr.AppError
		claims        models.Claims
	}{
		{
			name: "successfully updated task",
			setupTaskRepo: func(mtr *mocks.MockTaskRepo) {
				mtr.EXPECT().UpdateTask(int64(1234), models.Task{
					Title:  "title",
					Desc:   "desc",
					Status: 0,
					UserID: 1234,
				}).Return(nil)
			},
			id: "1234",
			taskReq: models.TaskRequestDto{
				Title:  "title",
				Desc:   "desc",
				Status: "Pending",
			},
			appErr: nil,
			claims: models.Claims{
				ID:   1234,
				Role: "",
			},
		},
		{
			name: "successfully created task with only title changed",
			setupTaskRepo: func(mtr *mocks.MockTaskRepo) {
				mtr.EXPECT().UpdateTask(int64(1234), models.Task{
					Title:  "title",
					Status: -1,
					UserID: 1234,
				}).Return(nil)
			},
			id: "1234",
			taskReq: models.TaskRequestDto{
				Title: "title",
			},
			appErr: nil,
			claims: models.Claims{
				ID:   1234,
				Role: "",
			},
		},
		{
			name:          "failed to update task because of invalid id",
			setupTaskRepo: func(mtr *mocks.MockTaskRepo) {},
			taskReq:       models.TaskRequestDto{},
			id:            "1234r",
			appErr: &errr.AppError{
				Message: "Invalid task id",
				Code:    http.StatusBadRequest,
			},
			claims: models.Claims{
				ID:   1234,
				Role: "",
			},
		},
		{
			name:          "failed to update task beacause to empty task fields",
			setupTaskRepo: func(mtr *mocks.MockTaskRepo) {},
			id:            "1234",
			taskReq:       models.TaskRequestDto{},
			appErr: &errr.AppError{
				Message: "Invalid task format",
				Code:    http.StatusBadRequest,
			},
			claims: models.Claims{
				ID:   1234,
				Role: "",
			},
		},
		{
			name: "task repo failed to update task",
			setupTaskRepo: func(mtr *mocks.MockTaskRepo) {
				mtr.EXPECT().UpdateTask(int64(1234), models.Task{
					Title:  "title",
					Desc:   "desc",
					Status: 0,
					UserID: 1234,
				}).Return(&errr.AppError{
					Code:    0,
					Message: "error message from task repo",
				})
			},
			id: "1234",
			taskReq: models.TaskRequestDto{
				Title:  "title",
				Desc:   "desc",
				Status: "Pending",
			},
			appErr: &errr.AppError{
				Code:    0,
				Message: "error message from task repo",
			},
			claims: models.Claims{
				ID:   1234,
				Role: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mtr := mocks.NewMockTaskRepo(ctrl)
			tt.setupTaskRepo(mtr)

			ts := NewTaskService(mtr)
			got := ts.UpdateTask(tt.id, tt.taskReq, tt.claims)

			if tt.appErr == nil && tt.appErr != got {
				t.Errorf("UpdateTask() failed, got err: %v.", got)
				return
			}
			if tt.appErr != nil && got == nil {
				t.Errorf("UpdateTask() successed unexpectedly, wanted err: %v.", tt.appErr)
				return
			}
			if tt.appErr != nil && *tt.appErr != *got {
				t.Errorf("UpdateTask() = %v, want %v", got, tt.appErr)
			}
		})
	}
}

func Test_taskService_DeleteTask(t *testing.T) {
	tests := []struct {
		name          string
		setupTaskRepo func(mtr *mocks.MockTaskRepo)
		id            string
		appErr        *errr.AppError
		claims        models.Claims
	}{
		{
			name: "successfully deleted task",
			setupTaskRepo: func(mtr *mocks.MockTaskRepo) {
				mtr.EXPECT().DeleteTask(int64(1234), int64(4321)).Return(nil)
			},
			id:     "1234",
			appErr: nil,
			claims: models.Claims{
				ID:   4321,
				Role: "",
			},
		},
		{
			name:          "failed to delete task because of invalid id",
			setupTaskRepo: func(mtr *mocks.MockTaskRepo) {},
			id:            "1234r",
			appErr: &errr.AppError{
				Message: "invalid id",
				Code:    http.StatusBadRequest,
			},
			claims: models.Claims{
				ID:   1234,
				Role: "",
			},
		},
		{
			name: "task repo failed to delete task",
			setupTaskRepo: func(mtr *mocks.MockTaskRepo) {
				mtr.EXPECT().DeleteTask(int64(1234), int64(1234)).Return(&errr.AppError{
					Code:    0,
					Message: "error message from task repo",
				})
			},
			id: "1234",
			appErr: &errr.AppError{
				Code:    0,
				Message: "error message from task repo",
			},
			claims: models.Claims{
				ID:   1234,
				Role: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mtr := mocks.NewMockTaskRepo(ctrl)
			tt.setupTaskRepo(mtr)
			ts := NewTaskService(mtr)

			got := ts.DeleteTask(tt.id, tt.claims)
			if tt.appErr == nil && tt.appErr != got {
				t.Errorf("DeleteTask() failed, got err: %v.", got)
				return
			}
			if tt.appErr != nil && got == nil {
				t.Errorf("DeleteTask() successed unexpectedly, wanted err: %v.", tt.appErr)
				return
			}
			if tt.appErr != nil && *tt.appErr != *got {
				t.Errorf("DeleteTask() = %v, want %v", got, tt.appErr)
			}
		})
	}
}

func Test_taskService_GetTasks(t *testing.T) {
	tests := []struct {
		name          string
		setupTaskRepo func(mtr *mocks.MockTaskRepo)
		appErr        *errr.AppError
		want          []models.TaskResponseDto
		claims        models.Claims
	}{
		{
			name: "successfully got task",
			setupTaskRepo: func(mtr *mocks.MockTaskRepo) {
				mtr.EXPECT().GetTasks(int64(1234)).Return([]models.Task{
					{
						ID:     1234,
						Title:  "my title",
						Desc:   "my string",
						Status: 0,
					},
					{
						ID:     1235,
						Title:  "my title",
						Desc:   "my string",
						Status: 1,
					},
				}, nil)
			},
			appErr: nil,
			want: []models.TaskResponseDto{
				{
					ID:     "1234",
					Title:  "my title",
					Desc:   "my string",
					Status: "Pending",
				},
				{
					ID:     "1235",
					Title:  "my title",
					Desc:   "my string",
					Status: "Done",
				},
			},
			claims: models.Claims{
				ID:   1234,
				Role: "",
			},
		},
		{
			name: "task repo failed to get task",
			setupTaskRepo: func(mtr *mocks.MockTaskRepo) {
				mtr.EXPECT().GetTasks(int64(1234)).Return(nil, &errr.AppError{
					Code:    0,
					Message: "error message from task repo",
				})
			},
			appErr: &errr.AppError{
				Code:    0,
				Message: "error message from task repo",
			},
			claims: models.Claims{
				ID:   1234,
				Role: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mtr := mocks.NewMockTaskRepo(ctrl)
			tt.setupTaskRepo(mtr)
			ts := NewTaskService(mtr)

			got, err := ts.GetTasks(tt.claims)

			if tt.appErr == nil && tt.appErr != err {
				t.Errorf("DeleteTask() failed, got err: %v.", got)
				return
			}
			if tt.appErr == nil && !slices.Equal(got, tt.want) {
				t.Errorf("wanted output: %v, got: %v", tt.want, got)
			}
			if tt.appErr != nil && err == nil {
				t.Errorf("DeleteTask() successed unexpectedly, wanted err: %v.", tt.appErr)
				return
			}
			if tt.appErr != nil && *tt.appErr != *err {
				t.Errorf("DeleteTask() = %v, want %v", got, tt.appErr)
			}
		})
	}
}

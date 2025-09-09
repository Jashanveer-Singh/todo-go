package services

import (
	"errors"
	"net/http"
	"testing"

	"github.com/Jashanveer-Singh/todo-go/internal/errr"
	"github.com/Jashanveer-Singh/todo-go/internal/models"
	"github.com/Jashanveer-Singh/todo-go/test/mocks"
	"github.com/golang/mock/gomock"
)

func Test_userService_CreateUser(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		userReq models.UserRequestDto
		// Named input parameters for receiver constructor.
		setupUserRepo       func(mur *mocks.MockUserRepo)
		setupPasswordHasher func(mph *mocks.MockPasswordHasher)
		// Named input parameters for target function.
		wantAppErr *errr.AppError
	}{
		{
			name: "empty user",
			userReq: models.UserRequestDto{
				Username: "",
				Password: "",
			},
			setupUserRepo:       func(mur *mocks.MockUserRepo) {},
			setupPasswordHasher: func(mph *mocks.MockPasswordHasher) {},
			wantAppErr: &errr.AppError{
				Code:    http.StatusBadRequest,
				Message: "Invalid user data",
			},
		},
		{
			name: "password too short",
			userReq: models.UserRequestDto{
				Username: "user",
				Password: "pass",
			},
			setupUserRepo:       func(mur *mocks.MockUserRepo) {},
			setupPasswordHasher: func(mph *mocks.MockPasswordHasher) {},
			wantAppErr: &errr.AppError{
				Code:    http.StatusBadRequest,
				Message: "Invalid user data",
			},
		},
		{
			name: "password hasher failed to create hash",
			userReq: models.UserRequestDto{
				Username: "user",
				Password: "password",
			},
			setupUserRepo: func(mur *mocks.MockUserRepo) {},
			setupPasswordHasher: func(mph *mocks.MockPasswordHasher) {
				mph.EXPECT().
					Hash("password").
					Return("", errors.New("error message from password hasher"))
			},
			wantAppErr: &errr.AppError{
				Code:    http.StatusInternalServerError,
				Message: "error message from password hasher",
			},
		},
		{
			name: "user repo failed to create user",
			userReq: models.UserRequestDto{
				Username: "user",
				Password: "password",
			},
			setupUserRepo: func(mur *mocks.MockUserRepo) {
				mur.EXPECT().CreateUser(models.User{
					ID:       0,
					Username: "user",
					Password: "hashed password",
				}).Return(&errr.AppError{
					Code:    0,
					Message: "error message from user repo",
				})
			},
			setupPasswordHasher: func(mph *mocks.MockPasswordHasher) {
				mph.EXPECT().Hash("password").Return("hashed password", nil)
			},
			wantAppErr: &errr.AppError{
				Code:    0,
				Message: "error message from user repo",
			},
		},
		{
			name: "successfully created user",
			userReq: models.UserRequestDto{
				Username: "user",
				Password: "password",
			},
			setupUserRepo: func(mur *mocks.MockUserRepo) {
				mur.EXPECT().CreateUser(models.User{
					ID:       0,
					Username: "user",
					Password: "hashed password",
				}).Return(nil)
			},
			setupPasswordHasher: func(mph *mocks.MockPasswordHasher) {
				mph.EXPECT().Hash("password").Return("hashed password", nil)
			},
			wantAppErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepoCtrl := gomock.NewController(t)
			defer userRepoCtrl.Finish()
			userRepo := mocks.NewMockUserRepo(userRepoCtrl)
			tt.setupUserRepo(userRepo)

			passwordHasherCtrl := gomock.NewController(t)
			defer passwordHasherCtrl.Finish()
			passwordHasher := mocks.NewMockPasswordHasher(passwordHasherCtrl)
			tt.setupPasswordHasher(passwordHasher)

			as := NewUserService(userRepo, passwordHasher)
			gotAppErr := as.CreateUser(tt.userReq)
			// TODO: update the condition below to compare got with tt.want.
			if tt.wantAppErr == nil && gotAppErr != nil {
				t.Errorf("CreateUser() failed. got appErr: %v", gotAppErr)
				return
			}
			if tt.wantAppErr != nil && gotAppErr == nil {
				t.Errorf("CreateUser succeessed unexpectedly. wanted appErr: %v", tt.wantAppErr)
				return
			}

			if tt.wantAppErr != nil && *gotAppErr != *tt.wantAppErr {
				t.Errorf("wanted appErr: %v, got: %v", tt.wantAppErr, gotAppErr)
			}
		})
	}
}

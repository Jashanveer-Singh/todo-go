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

func Test_authService_Login(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		setupUserRepo       func(mur *mocks.MockUserRepo)
		setupTokenProvider  func(mtp *mocks.MockTokenProvider)
		setupPasswordHasher func(mph *mocks.MockPasswordHasher)
		// Named input parameters for target function.
		username   string
		password   string
		want       string
		wantAppErr *errr.AppError
	}{
		{
			name:     "empty username",
			username: "",
			setupUserRepo: func(mur *mocks.MockUserRepo) {
				mur.EXPECT().GetUserByUsername("").Return(models.User{}, &errr.AppError{
					Code:    0,
					Message: "error message from user repo",
				})
			},
			setupTokenProvider:  func(mtp *mocks.MockTokenProvider) {},
			setupPasswordHasher: func(mph *mocks.MockPasswordHasher) {},
			want:                "",
			wantAppErr: &errr.AppError{
				Code:    0,
				Message: "error message from user repo",
			},
		},
		{
			name:     "empty password",
			username: "user",
			password: "",
			setupUserRepo: func(mur *mocks.MockUserRepo) {
				mur.EXPECT().GetUserByUsername("user").Return(models.User{
					ID:       0,
					Username: "user",
					Password: "password",
				}, nil)
			},
			setupPasswordHasher: func(mph *mocks.MockPasswordHasher) {
				mph.EXPECT().
					CompareHash("password", "").
					Return(false, errors.New("error message from password hasher"))
			},
			setupTokenProvider: func(mtp *mocks.MockTokenProvider) {},
			want:               "",
			wantAppErr: &errr.AppError{
				Code:    http.StatusInternalServerError,
				Message: "error message from password hasher",
			},
		},
		{
			name:     "user not found",
			username: "user",
			setupUserRepo: func(mur *mocks.MockUserRepo) {
				mur.EXPECT().GetUserByUsername("user").Return(models.User{}, &errr.AppError{
					Code:    0,
					Message: "error message from user repo, user not found",
				})
			},
			setupPasswordHasher: func(mph *mocks.MockPasswordHasher) {},
			setupTokenProvider:  func(mtp *mocks.MockTokenProvider) {},
			want:                "",
			wantAppErr: &errr.AppError{
				Code:    0,
				Message: "error message from user repo, user not found",
			},
		},
		{
			name:     "invalid password",
			username: "user",
			password: "invalid password",
			setupUserRepo: func(mur *mocks.MockUserRepo) {
				mur.EXPECT().GetUserByUsername("user").Return(models.User{
					ID:       0,
					Username: "user",
					Password: "password",
				}, nil)
			},
			setupPasswordHasher: func(mph *mocks.MockPasswordHasher) {
				mph.EXPECT().
					CompareHash("password", "invalid password").
					Return(false, nil)
			},
			setupTokenProvider: func(mtp *mocks.MockTokenProvider) {},
			want:               "",
			wantAppErr: &errr.AppError{
				Code:    http.StatusUnauthorized,
				Message: "Invalid Username or password",
			},
		},
		{
			name:     "failed to generate token",
			username: "user",
			password: "password",
			setupUserRepo: func(mur *mocks.MockUserRepo) {
				mur.EXPECT().GetUserByUsername("user").Return(models.User{
					ID:       0,
					Username: "user",
					Password: "password",
				}, nil)
			},
			setupPasswordHasher: func(mph *mocks.MockPasswordHasher) {
				mph.EXPECT().
					CompareHash("password", "password").
					Return(true, nil)
			},
			setupTokenProvider: func(mtp *mocks.MockTokenProvider) {
				mtp.EXPECT().GenerateToken(models.Claims{ID: 0, Role: ""}).
					Return("", errors.New("error message from token provider"))
			},
			want: "",
			wantAppErr: &errr.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to create token",
			},
		},
		{
			name:     "successfully created token",
			username: "user",
			password: "password",
			setupUserRepo: func(mur *mocks.MockUserRepo) {
				mur.EXPECT().GetUserByUsername("user").Return(models.User{
					ID:       0,
					Username: "user",
					Password: "password",
				}, nil)
			},
			setupPasswordHasher: func(mph *mocks.MockPasswordHasher) {
				mph.EXPECT().
					CompareHash("password", "password").
					Return(true, nil)
			},
			setupTokenProvider: func(mtp *mocks.MockTokenProvider) {
				mtp.EXPECT().GenerateToken(models.Claims{ID: 0, Role: ""}).
					Return("token", nil)
			},
			want:       "token",
			wantAppErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepoCtrl := gomock.NewController(t)
			defer userRepoCtrl.Finish()
			userRepo := mocks.NewMockUserRepo(userRepoCtrl)
			tt.setupUserRepo(userRepo)

			tokenProviderCtrl := gomock.NewController(t)
			defer tokenProviderCtrl.Finish()
			tokenProvider := mocks.NewMockTokenProvider(tokenProviderCtrl)
			tt.setupTokenProvider(tokenProvider)

			passwordHasherCtrl := gomock.NewController(t)
			defer passwordHasherCtrl.Finish()
			passwordHasher := mocks.NewMockPasswordHasher(passwordHasherCtrl)
			tt.setupPasswordHasher(passwordHasher)

			as := NewAuthService(userRepo, tokenProvider, passwordHasher)
			got, gotAppErr := as.Login(tt.username, tt.password)
			if tt.wantAppErr == nil && gotAppErr != nil {
				t.Errorf("Login() failed, got err: %v", gotAppErr)
				return
			}
			if tt.wantAppErr == nil && gotAppErr == nil {
				if tt.want != got {
					t.Errorf("Login = %s, wanted: %s", got, tt.want)
				}
				return
			}
			if tt.wantAppErr != nil && gotAppErr == nil {
				t.Errorf("Login() successed unexpectedly, wanted err: %v", tt.wantAppErr)
				return
			}
			if tt.wantAppErr != nil && *tt.wantAppErr != *gotAppErr {
				t.Errorf("wanted err: %v, got %v", tt.wantAppErr, gotAppErr)
				return
			}
		})
	}
}

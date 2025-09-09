package file

import (
	"net/http"
	"os"
	"path"
	"slices"
	"strings"
	"testing"

	"github.com/Jashanveer-Singh/todo-go/internal/errr"
	"github.com/Jashanveer-Singh/todo-go/internal/models"
)

func getTempUsersPath(t *testing.T) string {
	return path.Join(t.TempDir(), "users.json")
}

func Test_userRepo_readUsersFromFile(t *testing.T) {
	tests := []struct {
		name      string
		fp        string
		setupFile func(fp string)
		want      []models.User
		wantErr   bool
		Err       string
	}{
		{
			name: "Invalid json",
			fp:   getTempUsersPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte("askd"), 0666)
			},
			want:    nil,
			wantErr: true,
			Err:     "unable to unmarshal/decode json",
		},
		{
			name: "Can't read from file",
			fp:   getTempUsersPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte(""), 0333)
			},
			want:    nil,
			wantErr: true,
			Err:     "unable to read Tasks from file",
		},
		{
			name: "successfully got tasks",
			fp:   getTempUsersPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte(`[{
					"id": 1234,
					"username": "user",
					"password": "password"
					}]`), 0644)
			},
			want: []models.User{
				{
					ID:       1234,
					Username: "user",
					Password: "password",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupFile(tt.fp)
			tr := NewUserRepo(tt.fp)
			got, gotErr := tr.readUsersFromFile()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("readUsersFromFile() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("readUsersFromFile() succeeded unexpectedly, got: ", got)
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("readUsersFromFile() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_userRepo_writeUsersToFile(t *testing.T) {
	tests := []struct {
		name      string
		fp        string
		setupFile func(fp string)
		user      []models.User
		wantErr   bool
		err       string
	}{
		{
			name: "empty users",
			fp:   getTempUsersPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte(""), 0333)
			},
			user:    []models.User{},
			wantErr: false,
		},
		{
			name: "unable to write user to file",
			fp:   getTempUsersPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte(""), 0644)
				os.Chmod(fp, 0444)
			},
			user: []models.User{
				{
					ID:       0,
					Username: "user",
					Password: "password",
				},
			},
			wantErr: true,
			err:     "failed to write users to file",
		},
		{
			name: "successfull write user",
			fp:   getTempUsersPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte(""), 0333)
			},
			user: []models.User{
				{
					ID:       0,
					Username: "user1",
					Password: "user1 password",
				},
				{
					ID:       1,
					Username: "user2",
					Password: "user2 password",
				},
			},
			wantErr: false,
		},
		{
			name: "incomplete user fields",
			fp:   getTempUsersPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte(""), 0333)
			},
			user: []models.User{
				{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupFile(tt.fp)
			tr := NewUserRepo(tt.fp)
			gotErr := tr.writeUsersToFile(tt.user)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("writeUsersToFile() failed: %v", gotErr)
				} else if !strings.Contains(gotErr.Error(), tt.err) {
					t.Errorf("got wrong error for writeUsersToFile(), wanted %v, got %v", tt.err, gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("writeUsersToFile() succeeded unexpectedly")
			}
		})
	}
}

func Test_userRepo_GetUserByUsername(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		fp        string
		setupFile func(fp string)
		// Named input parameters for target function.
		username string
		want     models.User
		appError *errr.AppError
	}{
		{
			name: "read to file failed",
			fp:   getTempUsersPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte("asdfs"), 0666)
			},
			username: "user",
			appError: &errr.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Unable to save user due to internal server error",
			},
		},
		{
			name: "user not found",
			fp:   getTempTasksPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte(`[{
					"id": 1234,
					"username": "user",
					"password": "password"
					}]`), 0666)
			},
			username: "otheruser",
			appError: &errr.AppError{
				Code:    http.StatusNotFound,
				Message: "User not Found",
			},
		},
		{
			name: "successfully got user",
			fp:   getTempUsersPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte(`[{
					"id": 1234,
					"username": "user",
					"password": "password"
					}]`), 0666)
			},
			username: "user",
			want: models.User{
				ID:       1234,
				Username: "user",
				Password: "password",
			},
			appError: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupFile(tt.fp)
			ur := NewUserRepo(tt.fp)
			gotUser, gotAppErr := ur.GetUserByUsername(tt.username)
			// TODO: update the condition below to compare got with tt.want.
			if tt.appError == nil && gotAppErr != nil {
				t.Errorf("GetUserByUsername() failed. wanted app err: %v", tt.appError)
				return
			}
			if tt.appError == nil && gotAppErr == nil {
				if tt.want != gotUser {
					t.Errorf("want user: %v, got: %v", tt.want, gotUser)
				}
				return
			}

			if gotAppErr == nil {
				t.Errorf(
					"GetUserByUsername() successed unexpectedly. wanted app err = %v",
					tt.appError,
				)
				return
			}
			if *tt.appError != *gotAppErr {
				t.Errorf("wanted app err: %v, got: %v.", tt.appError, gotAppErr)
			}
		})
	}
}

func Test_userRepo_CreateUser(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		fp        string
		setupFile func(fp string)
		// Named input parameters for target function.
		user       models.User
		wantAppErr *errr.AppError
	}{
		{
			name: "read to file failed",
			fp:   getTempUsersPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte("asdfs"), 0666)
			},
			wantAppErr: &errr.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Unable to save user due to internal server error",
			},
		},
		{
			name: "user already exists",
			fp:   getTempUsersPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte(`[{
					"id": 12356,
					"username": "user",
					"password": "user password"
					}]`), 0666)
			},
			user: models.User{
				ID:       123,
				Username: "user",
			},
			wantAppErr: &errr.AppError{
				Code:    http.StatusConflict,
				Message: "user already exists",
			},
		},
		{
			name: "successfully created user",
			fp:   getTempUsersPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte(`[{
					"id": 12356,
					"username": "user",
					"password": "user password"
					}]`), 0666)
			},
			user: models.User{
				Username: "other user",
				Password: "password",
			},
			wantAppErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupFile(tt.fp)
			ur := NewUserRepo(tt.fp)
			gotAppErr := ur.CreateUser(tt.user)
			if tt.wantAppErr == nil && gotAppErr != nil {
				t.Errorf("CreateUser() failed. wanted app err: %v", tt.wantAppErr)
				return
			}
			if tt.wantAppErr != nil && gotAppErr == nil {
				t.Errorf(
					"CreateUser() successed unexpectedly, wanted app err: %v",
					tt.wantAppErr,
				)
				return
			}

			if tt.wantAppErr != nil && gotAppErr != nil {
				if *tt.wantAppErr != *gotAppErr {
					t.Errorf("want app err: %v, got %v", tt.wantAppErr, gotAppErr)
				}
			}
		})
	}
}

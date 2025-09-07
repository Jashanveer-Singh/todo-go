package file

import (
	"bytes"
	"net/http"
	"os"
	"os/exec"
	"path"
	"slices"
	"strings"
	"testing"

	"github.com/Jashanveer-Singh/todo-go/internal/errr"
	"github.com/Jashanveer-Singh/todo-go/internal/models"
)

func getTempTasksPath(t *testing.T) string {
	return path.Join(t.TempDir(), "tasks.json")
}

func TestNewTaskRepo_ValidFile(t *testing.T) {
	dir := t.TempDir()
	fp := path.Join(dir, "tasks.json")
	os.WriteFile(fp, []byte(""), 0644)

	got := NewTaskRepo(fp)

	if got.fp != fp {
		t.Errorf("wanted file pointer(fp): %s, got: %s.", fp, got.fp)
	}
}

func TestNewtaskRepo_FileDoesNotExist(t *testing.T) {
	dir := t.TempDir()
	fp := path.Join(dir, "tasks.json")

	if os.Getenv("BE_CRASHER") == "1" {
		NewTaskRepo(fp)
	}

	var stderr bytes.Buffer

	cmd := exec.Command(os.Args[0], "-test.run=TestNewtaskRepo_FileDoesNotExist")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	cmd.Stderr = &stderr
	err := cmd.Run()

	if e, ok := err.(*exec.ExitError); !ok || e.Success() {
		t.Errorf("wanted exit Error, got err: %v", e)
	}
	stderrOutput := stderr.String()

	if !strings.Contains(stderrOutput, "does not exist") {
		t.Errorf("wanted does not exist error, got %s.", stderrOutput)
	}
}

// func TestNewtaskRepo_NotEnoughPermissionsForFile(t *testing.T) {
// 	if os.Getenv("BE_CRASHER") == "1" {
// 		dir := t.TempDir()
// 		fp := path.Join(dir, "tasks.json")
// 		err := os.WriteFile(fp, []byte("asdk"), 0600)
// 		if err != nil {
// 			os.Stderr.Write([]byte("failed to create file"))
// 		}
// 		err = os.Chmod(fp, 0000)
// 		if err != nil {
// 			os.Stderr.Write([]byte("failed to set file permissions"))
// 		}
// 		NewTaskRepo(fp)
// 		return
// 	}
//
// 	var stderr bytes.Buffer
//
// 	cmd := exec.Command(os.Args[0], "-test.run=TestNewtaskRepo_NotEnoughPermissionsForFile")
// 	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
// 	cmd.Stderr = &stderr
// 	err := cmd.Run()
//
// 	if e, ok := err.(*exec.ExitError); !ok || e.Success() {
// 		t.Errorf("wanted exit Error, got err: %v", e)
// 	}
// 	stderrOutput := stderr.String()
//
// 	if !strings.Contains(stderrOutput, "not enough permissions for file") {
// 		t.Errorf("wanted not enough permissions error, got %s.", stderrOutput)
// 	}
// }

func Test_taskRepo_getTasks(t *testing.T) {
	tests := []struct {
		name      string
		fp        string
		setupFile func(fp string)
		want      []models.Task
		wantErr   bool
		Err       string
	}{
		{
			name: "Invalid json",
			fp:   getTempTasksPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte("askd"), 0666)
			},
			want:    nil,
			wantErr: true,
			Err:     "unable to unmarshal/decode json",
		},
		{
			name: "Can't read from file",
			fp:   getTempTasksPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte(""), 0333)
			},
			want:    nil,
			wantErr: true,
			Err:     "unable to read Tasks from file",
		},
		{
			name: "successfully got tasks",
			fp:   getTempTasksPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte(`[{
					"title": "any title",
					"desc": "any desc",
					"status": 0,
					"id": 12234
					}]`), 0644)
			},
			want: []models.Task{
				{
					ID:     12234,
					Status: 0,
					Desc:   "any desc",
					Title:  "any title",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupFile(tt.fp)
			tr := NewTaskRepo(tt.fp)
			got, gotErr := tr.getTasks()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("getTasks() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("getTasks() succeeded unexpectedly, got: ", got)
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("getTasks() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_taskRepo_write(t *testing.T) {
	tests := []struct {
		name      string
		fp        string
		setupFile func(fp string)
		tasks     []models.Task
		wantErr   bool
		err       string
	}{
		{
			name: "empty tasks",
			fp:   getTempTasksPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte(""), 0333)
			},
			tasks:   []models.Task{},
			wantErr: false,
		},
		{
			name: "unable to write to file",
			fp:   getTempTasksPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte(""), 0644)
				os.Chmod(fp, 0444)
			},
			tasks: []models.Task{
				{
					ID:     0,
					Title:  "some Title",
					Desc:   "some desc",
					Status: 0,
				},
			},
			wantErr: true,
			err:     "unable to write tasks to file",
		},
		{
			name: "successfull write task",
			fp:   getTempTasksPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte(""), 0333)
			},
			tasks: []models.Task{
				{
					ID:     0,
					Title:  "some Title",
					Desc:   "some desc",
					Status: 0,
				},
				{
					ID:     1,
					Title:  "some other title",
					Desc:   "some other desc",
					Status: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "incomplete task fields",
			fp:   getTempTasksPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte(""), 0333)
			},
			tasks: []models.Task{
				{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupFile(tt.fp)
			tr := NewTaskRepo(tt.fp)
			gotErr := tr.write(tt.tasks)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("write() failed: %v", gotErr)
				} else if !strings.Contains(gotErr.Error(), tt.err) {
					t.Errorf("got wrong error for write(), wanted %v, got %v", tt.err, gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("write() succeeded unexpectedly")
			}
		})
	}
}

func Test_taskRepo_SaveTask(t *testing.T) {
	tests := []struct {
		name      string
		fp        string
		setupFile func(fp string)
		task      models.Task
		wantErr   bool
	}{
		{
			name: "empty task",
			fp:   getTempTasksPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte(""), 0666)
			},
			task:    models.Task{},
			wantErr: false,
		},
		{
			name: "tasks read failure due to corrupted file",
			fp:   getTempTasksPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte("asdf"), 0666)
			},
			task: models.Task{
				ID:     0,
				Title:  "some Title",
				Desc:   "some desc",
				Status: 0,
			},
			wantErr: true,
		},
		{
			name: "task write failure",
			fp:   getTempTasksPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte(""), 0444)
			},
			task: models.Task{
				ID:     0,
				Title:  "some Title",
				Desc:   "some desc",
				Status: 0,
			},
			wantErr: true,
		},
		{
			name: "task written successfully",
			fp:   path.Join(t.TempDir(), "tasks.join"),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte(`[{
					"title": "any title",
					"desc": "any desc",
					"status": 0,
					"id": 12234
					}]`), 0666)
			},
			task: models.Task{
				ID:     0,
				Title:  "some Title",
				Desc:   "some desc",
				Status: 0,
			},
			wantErr: false,
		},
		{
			name: "writing first task",
			fp:   path.Join(t.TempDir(), "task.json"),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte(""), 0666)
			},
			task: models.Task{
				ID:     0,
				Title:  "some Title",
				Desc:   "some desc",
				Status: 0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupFile(tt.fp)
			tr := NewTaskRepo(tt.fp)
			gotErr := tr.SaveTask(tt.task)
			if tt.wantErr && gotErr == nil {
				t.Errorf("SaveTask() successed unexpectedly")
			}
			if !tt.wantErr && gotErr != nil {
				t.Errorf("SaveTask failed. got %v", gotErr)
			}
		})
	}
}

func Test_taskRepo_UpdateTask(t *testing.T) {
	tests := []struct {
		name       string
		fp         string
		setupFile  func(fp string)
		id         int64
		task       models.Task
		wantErr    bool
		errMessage string
	}{
		{
			name: "task not found",
			fp:   getTempTasksPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte(""), 0666)
			},
			id:         0,
			task:       models.Task{},
			wantErr:    true,
			errMessage: "no task found with id",
		},
		{
			name: "empty task",
			fp:   getTempTasksPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte(`[{
					"title": "any title",
					"desc": "any desc",
					"status": 0,
					"id": 12234
					}]`), 0666)
			},
			id:      12234,
			task:    models.Task{},
			wantErr: false,
			// errMessage: "no task found with id",
		},
		{
			name: "tasks read failure due to corrupted file",
			fp:   getTempTasksPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte(`adsf`), 0666)
			},
			id:         12234,
			task:       models.Task{},
			wantErr:    true,
			errMessage: "Unable to create task due to internal server error",
		},
		{
			name: "unable to write tasks",
			fp:   getTempTasksPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte(`adsf`), 0444)
			},
			id:         12234,
			task:       models.Task{},
			wantErr:    true,
			errMessage: "Unable to create task due to internal server error",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupFile(tt.fp)
			tr := NewTaskRepo(tt.fp)
			gotErr := tr.UpdateTask(tt.id, tt.task)
			// TODO: update the condition below to compare got with tt.want.
			if tt.wantErr && gotErr == nil {
				// t.Errorf("UpdateTask() = %v, want %v", got, tt.want)
				t.Errorf("UpdateTask() successed unexpectedly")
			}
			if !tt.wantErr && gotErr != nil {
				t.Errorf("UpdateTask() failed, got err %v", gotErr.Message)
			}
			if tt.wantErr && gotErr != nil {
				if gotErr.Message != tt.errMessage {
					t.Errorf("wanted error %v, got err %v", tt.errMessage, gotErr.Message)
				}
			}
		})
	}
}

func Test_taskRepo_DeleteTask(t *testing.T) {
	tests := []struct {
		name       string
		fp         string
		setupFile  func(fp string)
		id         int64
		wantErr    bool
		errMessage string
	}{
		{
			name: "task not found",
			fp:   getTempTasksPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte(""), 0666)
			},
			id:         0,
			wantErr:    true,
			errMessage: "no task found with id",
		},
		{
			name: "tasks read failure due to corrupted file",
			fp:   getTempTasksPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte(`adsf`), 0666)
			},
			id:         12234,
			wantErr:    true,
			errMessage: "Unable to create task due to internal server error",
		},
		{
			name: "unable to write tasks",
			fp:   getTempTasksPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte(`adsf`), 0444)
			},
			id:         12234,
			wantErr:    true,
			errMessage: "Unable to create task due to internal server error",
		},
		{
			name: "task deleted successfully",
			fp:   getTempTasksPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte(`[{
					"title": "any title",
					"desc": "any desc",
					"status": 0,
					"id": 12234
					}]`), 0666)
			},
			id:      12234,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupFile(tt.fp)
			tr := NewTaskRepo(tt.fp)
			gotErr := tr.DeleteTask(tt.id)
			if tt.wantErr && gotErr == nil {
				t.Errorf("UpdateTask() successed unexpectedly")
			}
			if !tt.wantErr && gotErr != nil {
				t.Errorf("UpdateTask() failed, got err %v", gotErr.Message)
			}
			if tt.wantErr && gotErr != nil {
				if gotErr.Message != tt.errMessage {
					t.Errorf("wanted error %v, got err %v", tt.errMessage, gotErr.Message)
				}
			}
		})
	}
}

func Test_taskRepo_GetTasks(t *testing.T) {
	tests := []struct {
		name      string
		fp        string
		setupFile func(fp string)
		want      []models.Task
		err       *errr.AppError
		wantErr   bool
	}{
		{
			name: "tasks read failure due to corrupted file",
			fp:   getTempTasksPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte("asdfaf"), 0666)
			},
			wantErr: true,
			err: &errr.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Unable to create task due to internal server error",
			},
		},
		{
			name: "tasks read successfully",
			fp:   getTempTasksPath(t),
			setupFile: func(fp string) {
				os.WriteFile(fp, []byte(`[{
					"title": "any title",
					"desc": "any desc",
					"status": 0,
					"id": 12234
					}]`), 0666)
			},
			want: []models.Task{
				{
					ID:     12234,
					Status: 0,
					Desc:   "any desc",
					Title:  "any title",
				},
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupFile(tt.fp)
			tr := NewTaskRepo(tt.fp)
			got, err := tr.GetTasks()
			// TODO: update the condition below to compare got with tt.want.
			if tt.wantErr && err == nil {
				t.Errorf("GetTasks successed unexpectedly")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("GetTasks() Failed, got err %v", err)
			}
			if !tt.wantErr && err == nil {
				if !slices.Equal(got, tt.want) {
					t.Errorf("Wanted %v, got %v", tt.want, got)
				}
			}
		})
	}
}

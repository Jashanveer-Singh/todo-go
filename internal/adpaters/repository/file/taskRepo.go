package file

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Jashanveer-Singh/todo-go/internal/errr"
	"github.com/Jashanveer-Singh/todo-go/internal/models"
)

func NewTaskRepo(fp string) *taskRepo {
	_, err := os.Stat(fp)
	if os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "file: %s does not exist\n", fp)
	}
	if os.IsPermission(err) {
		fmt.Fprintf(os.Stderr, "not enough permissions for file: %s\n", fp)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "can't use the file: %s\n%s", fp, err.Error())
		os.Exit(1)
	}

	return &taskRepo{
		mu: sync.RWMutex{},
		fp: fp,
	}
}

type taskRepo struct {
	mu sync.RWMutex
	fp string
}

func (tr *taskRepo) getTasks() ([]models.Task, error) {
	tasks := make([]models.Task, 0)

	taskjson, err := os.ReadFile(tr.fp)
	if err != nil {
		return nil, fmt.Errorf("unable to read Tasks from file.\n%s", err.Error())
	}
	if len(taskjson) != 0 {

		err = json.Unmarshal(taskjson, &tasks)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal/decode json.\n%s", err.Error())
		}
	}

	return tasks, nil
}

func (tr *taskRepo) write(tasks []models.Task) error {
	taskjson, _ := json.Marshal(tasks)

	err := os.WriteFile(tr.fp, taskjson, 0644)
	if err != nil {
		return fmt.Errorf("unable to write tasks to file.\n%s", err.Error())
	}

	return nil
}

func (tr *taskRepo) SaveTask(task models.Task) *errr.AppError {
	task.ID = time.Now().Unix()
	tr.mu.Lock()
	defer tr.mu.Unlock()
	tasks, err := tr.getTasks()
	if err != nil {
		return errr.NewUnexpectedError("Unable to save task due to internal server error")
	}

	tasks = append(tasks, task)

	err = tr.write(tasks)
	if err != nil {
		return errr.NewUnexpectedError("Unable to save task due to internal server error")
	}

	return nil
}

func (tr *taskRepo) UpdateTask(id int64, task models.Task) *errr.AppError {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	tasks, err := tr.getTasks()
	if err != nil {
		return errr.NewUnexpectedError("Unable to create task due to internal server error")
	}

	notFound := true

	for i := range tasks {
		if tasks[i].ID == id {
			notFound = false
			if tasks[i].UserID != task.UserID {
				return errr.NewUnauthorizedError("Unauthorized to update task")
			}
			if len(task.Title) != 0 {
				tasks[i].Title = task.Title
			}
			if len(task.Desc) != 0 {
				tasks[i].Desc = task.Desc
			}
			if task.IsValidStatus() {
				tasks[i].Status = task.Status
			}
			break
		}
	}

	if notFound {
		return errr.NewNotFoundError("no task found with id")
	}

	err = tr.write(tasks)
	if err != nil {
		return errr.NewUnexpectedError("Unable to create task due to internal server error")
	}

	return nil
}

func (tr *taskRepo) DeleteTask(id int64, userID int64) *errr.AppError {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	tasks, err := tr.getTasks()
	if err != nil {
		return errr.NewUnexpectedError("Unable to create task due to internal server error")
	}

	notFound := true

	for i := range tasks {
		if tasks[i].ID == id {
			notFound = false
			if tasks[i].UserID != userID {
				return errr.NewUnauthorizedError("Unauthorized to delete task")
			}
			tasks = append(tasks[:i], tasks[i+1:]...)
			break
		}
	}

	if notFound {
		return errr.NewUnexpectedError("no task found with id")
	}

	err = tr.write(tasks)
	if err != nil {
		return errr.NewUnexpectedError("Unable to create task due to internal server error")
	}

	return nil
}

func (tr *taskRepo) GetTasks(userID int64) ([]models.Task, *errr.AppError) {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	tasks, err := tr.getTasks()
	if err != nil {
		return nil, errr.NewUnexpectedError("Unable to create task due to internal server error")
	}
	filteredTasks := []models.Task{}
	for _, task := range tasks {
		if task.UserID == userID {
			filteredTasks = append(filteredTasks, task)
		}
	}

	return filteredTasks, nil
}

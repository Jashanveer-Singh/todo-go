package file

import (
	"encoding/json"
	"errors"
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
		fmt.Fprintf(os.Stderr, "file: %s does not exist", fp)
	}
	if os.IsPermission(err) {
		fmt.Fprintf(os.Stderr, "not enough permissions for file: %s", fp)
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

func newRepoErr(err error) error {
	return errors.New(fmt.Sprint("Repo Error: ", err.Error()))
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
	taskjson, err := json.Marshal(tasks)
	if err != nil {
		return fmt.Errorf("unable to marshal tasks to json.\n%s", err.Error())
	}

	err = os.WriteFile(tr.fp, taskjson, 0644)
	if err != nil {
		return fmt.Errorf("unable to write tasks to file.\n%s", err.Error())
	}

	return nil
}

func (tr *taskRepo) SaveTask(task models.Task) *errr.AppError {
	task.Status = 0
	task.ID = time.Now().Unix()
	tr.mu.Lock()
	defer tr.mu.Unlock()
	tasks, err := tr.getTasks()
	if err != nil {
		return errr.NewUnexpectedError("Unable to create task due to internal server error")
	}

	tasks = append(tasks, task)

	err = tr.write(tasks)
	if err != nil {
		return errr.NewUnexpectedError("Unable to create task due to internal server error")
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

func (tr *taskRepo) DeleteTask(id int64) *errr.AppError {
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
			tasks = append(tasks[:i], tasks[i+1:]...)
			break
		}
	}

	if notFound {
		return errr.NewUnexpectedError("No task found with id")
	}

	err = tr.write(tasks)
	if err != nil {
		return errr.NewUnexpectedError("Unable to create task due to internal server error")
	}

	return nil
}

func (tr *taskRepo) GetTasks() ([]models.Task, *errr.AppError) {
	tr.mu.RLock()
	tr.mu.RUnlock()

	tasks, err := tr.getTasks()
	if err != nil {
		return nil, errr.NewUnexpectedError("Unable to create task due to internal server error")
	}

	return tasks, nil
}

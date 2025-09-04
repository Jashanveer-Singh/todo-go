package file

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/Jashanveer-Singh/todo-go/internal/handlers/domain"
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

func (tr *taskRepo) getTasks() ([]domain.Task, error) {
	tr.mu.RLock()
	defer tr.mu.Unlock()

	tasks := make([]domain.Task, 0)

	taskjson, err := os.ReadFile(tr.fp)
	if err != nil {
		return nil, fmt.Errorf("unable to read Tasks from file.\n%s", err.Error())
	}
	err = json.Unmarshal(taskjson, &tasks)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal/decode json.\n%s", err.Error())
	}

	return tasks, nil
}

func (tr *taskRepo) write(tasks []domain.Task) error {
	taskjson, err := json.Marshal(tasks)
	if err != nil {
		return fmt.Errorf("unable to marshal tasks to json.\n%s", err.Error())
	}

	tr.mu.Lock()
	defer tr.mu.Unlock()

	err = os.WriteFile(tr.fp, taskjson, 0644)
	if err != nil {
		return fmt.Errorf("unable to write tasks to file.\n%s", err.Error())
	}

	return nil
}

func (tr *taskRepo) CreateTask(task domain.Task) error {
	tasks, err := tr.getTasks()
	if err != nil {
		return newRepoErr(err)
	}

	tasks = append(tasks, task)

	err = tr.write(tasks)
	if err != nil {
		return newRepoErr(err)
	}

	return nil
}

func (tr *taskRepo) UpdateTask(id string, task domain.Task) error {
	tasks, err := tr.getTasks()
	if err != nil {
		return newRepoErr(err)
	}

	notFound := true

	for i := range tasks {
		if tasks[i].ID == id {
			notFound = false
			tasks[i] = task
			break
		}
	}

	if notFound {
		return newRepoErr(errors.New("no task found with id"))
	}

	err = tr.write(tasks)
	if err != nil {
		return newRepoErr(err)
	}

	return nil
}

func (tr *taskRepo) DeleteTask(id string) error {
	tasks, err := tr.getTasks()
	if err != nil {
		return newRepoErr(err)
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
		return newRepoErr(errors.New("no task found with id"))
	}

	err = tr.write(tasks)
	if err != nil {
		return newRepoErr(err)
	}

	return nil
}

func (tr *taskRepo) GetTasks() ([]domain.Task, error) { return tr.getTasks() }

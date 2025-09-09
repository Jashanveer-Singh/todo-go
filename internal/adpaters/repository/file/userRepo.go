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

func NewUserRepo(fp string) *userRepo {
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

	return &userRepo{
		mu: sync.RWMutex{},
		fp: fp,
	}
}

type userRepo struct {
	mu sync.RWMutex
	fp string
}

func (ur *userRepo) readUsersFromFile() ([]models.User, error) {
	users := []models.User{}

	userjson, err := os.ReadFile(ur.fp)
	if err != nil {
		return nil, fmt.Errorf("unable to read Tasks from file.\n%s", err.Error())
	}
	if len(userjson) != 0 {

		err = json.Unmarshal(userjson, &users)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal/decode json.\n%s", err.Error())
		}
	}

	return users, nil
}

func (ur *userRepo) writeUsersToFile(users []models.User) error {
	userjson, _ := json.Marshal(users)

	err := os.WriteFile(ur.fp, []byte(userjson), 0666)
	if err != nil {
		return fmt.Errorf("failed to write users to file.\n%s", err.Error())
	}

	return nil
}

func (ur *userRepo) GetUserByUsername(username string) (models.User, *errr.AppError) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()

	users, err := ur.readUsersFromFile()
	if err != nil {
		return models.User{}, errr.NewUnexpectedError(
			"Unable to save user due to internal server error",
		)
	}

	for i := range users {
		if users[i].Username == username {
			return users[i], nil
		}
	}

	return models.User{}, errr.NewNotFoundError("User not Found")
}

func (ur *userRepo) CreateUser(user models.User) *errr.AppError {
	ur.mu.RLock()
	defer ur.mu.RUnlock()

	users, err := ur.readUsersFromFile()
	if err != nil {
		return errr.NewUnexpectedError("Unable to save user due to internal server error")
	}

	for i := range users {
		if users[i].Username == user.Username {
			return errr.NewDuplicateError("user already exists")
		}
	}

	user.ID = time.Now().Unix()
	users = append(users, user)

	err = ur.writeUsersToFile(users)
	if err != nil {
		return errr.NewUnexpectedError("Unable to save user due to internal server error")
	}

	return nil
}

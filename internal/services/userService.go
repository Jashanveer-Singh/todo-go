package services

import (
	"strconv"

	"github.com/Jashanveer-Singh/todo-go/internal/errr"
	"github.com/Jashanveer-Singh/todo-go/internal/models"
	"github.com/Jashanveer-Singh/todo-go/internal/ports"
)

func NewUserService(ur ports.UserRepo) *userService {
	return &userService{
		ur: ur,
	}
}

type userService struct {
	ur ports.UserRepo
}

func (us *userService) CreateUser(userReq models.UserRequestDto) *errr.AppError {
	user := userReq.ToUser()

	if !user.IsValidUser() {
		return errr.NewBadRequestError("Invalid user data")
	}

	appErr := us.ur.CreateUser(user)
	if appErr != nil {
		return appErr
	}

	return nil
}

func (us *userService) Login(username, password string) (string, *errr.AppError) {
	user, appErr := us.ur.GetUserByUsername(username)
	if appErr != nil {
		return "", appErr
	}

	if user.Password == password {
		idString := strconv.FormatInt(user.ID, 10)
		return idString, nil
	}

	return "", errr.NewUnauthenticatedError("Invalid Username or password")
}

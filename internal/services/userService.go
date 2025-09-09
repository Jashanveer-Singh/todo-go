package services

import (
	"github.com/Jashanveer-Singh/todo-go/internal/errr"
	"github.com/Jashanveer-Singh/todo-go/internal/models"
	"github.com/Jashanveer-Singh/todo-go/internal/ports"
)

func NewUserService(userRepo ports.UserRepo, passwordHasher ports.PasswordHasher) *userService {
	return &userService{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
	}
}

type userService struct {
	userRepo       ports.UserRepo
	passwordHasher ports.PasswordHasher
}

func (as *userService) CreateUser(userReq models.UserRequestDto) *errr.AppError {
	user := userReq.ToUser()

	if !user.IsValidUser() {
		return errr.NewBadRequestError("Invalid user data")
	}
	hash, err := as.passwordHasher.Hash(user.Password)
	if err != nil {
		return errr.NewUnexpectedError(err.Error())
	}
	user.Password = hash

	appErr := as.userRepo.CreateUser(user)
	if appErr != nil {
		return appErr
	}

	return nil
}

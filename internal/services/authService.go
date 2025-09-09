package services

import (
	"github.com/Jashanveer-Singh/todo-go/internal/errr"
	"github.com/Jashanveer-Singh/todo-go/internal/models"
	"github.com/Jashanveer-Singh/todo-go/internal/ports"
)

func NewAuthService(
	userRepo ports.UserRepo,
	tokenProvider ports.TokenProvider,
	passwordHasher ports.PasswordHasher,
) *authService {
	return &authService{
		userRepo:       userRepo,
		tokenProvider:  tokenProvider,
		passwordHasher: passwordHasher,
	}
}

type authService struct {
	userRepo       ports.UserRepo
	tokenProvider  ports.TokenProvider
	passwordHasher ports.PasswordHasher
}

func (as *authService) Login(username, password string) (string, *errr.AppError) {
	user, appErr := as.userRepo.GetUserByUsername(username)
	if appErr != nil {
		return "", appErr
	}

	match, err := as.passwordHasher.CompareHash(user.Password, password)
	if err != nil {
		return "", errr.NewUnexpectedError(err.Error())
	}
	if !match {
		return "", errr.NewUnauthenticatedError("Invalid Username or password")
	}

	claims := models.Claims{
		ID:   user.ID,
		Role: "",
	}

	token, err := as.tokenProvider.GenerateToken(claims)
	if err != nil {
		return "", errr.NewUnexpectedError("Failed to create token")
	}
	return token, nil
}

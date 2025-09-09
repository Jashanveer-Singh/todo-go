package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/Jashanveer-Singh/todo-go/internal/adpaters/apis/http"
	"github.com/Jashanveer-Singh/todo-go/internal/adpaters/bcrypt"
	"github.com/Jashanveer-Singh/todo-go/internal/adpaters/jwttoken"
	"github.com/Jashanveer-Singh/todo-go/internal/adpaters/repository/file"
	"github.com/Jashanveer-Singh/todo-go/internal/services"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't get the current working directory\n")
		os.Exit(1)
	}
	dirPath := path.Join(cwd, "data")

	tasksFile := path.Join(dirPath, "tasks.json")
	usersFile := path.Join(dirPath, "users.json")

	taskRepo := file.NewTaskRepo(tasksFile)
	userRepo := file.NewUserRepo(usersFile)
	jwtTokenProvider := jwttoken.NewJWTTokenProvider(
		"my secret key",
		"issuer",
		"audience",
		time.Hour*24,
	)
	bcryptPasswordHasher := bcrypt.NewBcryptPasswordHasher(10)

	authService := services.NewAuthService(userRepo, jwtTokenProvider, bcryptPasswordHasher)
	taskService := services.NewTaskService(taskRepo)
	userService := services.NewUserService(userRepo, bcryptPasswordHasher)
	apiServer := http.NewHttpServer(taskService, userService, authService, jwtTokenProvider)

	log.Println("Starting Server at port:8080")
	apiServer.ListenAndServe(":8080")
}

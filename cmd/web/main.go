package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/Jashanveer-Singh/todo-go/internal/adpaters/apis/http"
	"github.com/Jashanveer-Singh/todo-go/internal/adpaters/repository/file"
	"github.com/Jashanveer-Singh/todo-go/internal/services"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"Can't get the home directory, the env variable for home directory may not be set\n",
		)
		os.Exit(1)
	}

	tasksFile := path.Join(homeDir, ".local", "todo", "tasks.json")
	usersFile := path.Join(homeDir, ".local", "todo", "users.json")

	taskRepo := file.NewTaskRepo(tasksFile)
	userRepo := file.NewUserRepo(usersFile)
	taskService := services.NewTaskService(taskRepo)
	userService := services.NewUserService(userRepo)
	apiServer := http.NewHttpServer(taskService, userService)

	log.Println("Starting Server at port:8080")
	apiServer.ListenAndServe(":8080")
}

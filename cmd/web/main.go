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

	fp := path.Join(homeDir, ".local", "todo", "tasks.json")

	repo := file.NewTaskRepo(fp)
	taskService := services.NewTaskService(repo)
	apiServer := http.NewHttpServer(taskService)

	log.Println("Starting Server at port:8080")
	apiServer.ListenAndServe(":8080")
}

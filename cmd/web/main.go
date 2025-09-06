package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/Jashanveer-Singh/todo-go/internal/adpaters/repository/file"
	"github.com/Jashanveer-Singh/todo-go/internal/handlers"
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
	handler := handlers.NewHandler(repo)
	router := handlers.NewRouter(handler)

	log.Println("Starting Server at port:8080")
	http.ListenAndServe(":8080", router)
}

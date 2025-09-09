package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var (
	totalTableWidth float32 = 100
	snoWidth                = 2
	titleWidth              = 0
	descWidth               = 0
	statusWidth             = 20
	tasks           []task
	username        string
	password        string
	token           string
)

func init() {
	titleWidth = int(totalTableWidth * 0.4)
	descWidth = int(totalTableWidth * 0.6)
}

type task struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Desc   string `json:"Desc"`
	Status string `json:"status"`
}

func pressEnterToContinue() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("\nPress Enter to Continue...")
	scanner.Scan()
}

func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func main() {
	input := -1

	clearScreen()
	fmt.Printf("1. Login\n2. Create User\nChoose: ")
	fmt.Scan(&input)

	switch input {
	case 1:
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Enter Username: ")
		scanner.Scan()
		username = scanner.Text()
		fmt.Print("Enter Password: ")
		scanner.Scan()
		password = scanner.Text()
		fmt.Printf("\n\n")
		handleLogin()
	case 2:
		handleCreateUser()
	default:
		printErrf("Invalid Option")
	}
}

func handleCreateUser() {
	clearScreen()
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter Username: ")
	scanner.Scan()
	username = scanner.Text()
	fmt.Print("Enter Password: ")
	scanner.Scan()
	password = scanner.Text()
	fmt.Printf("\n\n")

	requestBody := fmt.Sprintf(`{"username":"%s", "password": "%s"}`, username, password)

	response, err := http.Post(
		"http://localhost:8080/users",
		"application/json",
		strings.NewReader(requestBody),
	)
	if err != nil {
		printErrf("unexpected error while http request.\n%s\n", err.Error())
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		data, err := io.ReadAll(response.Body)
		if err != nil {
			printErrf("Failed to create user and failed to get error message")
			return
		}
		printErrf("Failed to create user. err: %s", string(data))
		return
	}
	fmt.Println("User created successfully")

	pressEnterToContinue()
	handleLogin()
}

func handleLogin() {
	requestBody := fmt.Sprintf(`{"username":"%s", "password": "%s"}`, username, password)

	response, err := http.Post(
		"http://localhost:8080/auth",
		"application/json",
		strings.NewReader(requestBody),
	)
	if err != nil {
		printErrf("unexpected error while http request.\n%s\n", err.Error())
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		data, err := io.ReadAll(response.Body)
		if err != nil {
			printErrf("Failed to get user token and error message")
			return
		}
		printErrf("Failed to get user token. err: %s", data)
		return
	}
	data, err := io.ReadAll(response.Body)
	if err != nil {
		printErrf("Failed to get user token. error in reading response body")
		return
	}
	token = "Bearer " + string(data)
	pressEnterToContinue()
	handlePostLogin()
}

func handlePostLogin() {
	for {
		clearScreen()
		input := -1
		fmt.Println("")
		fmt.Println("--------------Menu--------------")
		fmt.Printf("1. Add Task\n2. View Tasks\n3. Update Task\n4. Delete Task\n5. Exit\nChoose: ")
		fmt.Scan(&input)
		fmt.Printf("\n\n")

		switch input {
		case 1:
			handleCreateTask()
		case 2:
			handleShowTask()
		case 3:
			handleUpdateTask()
		case 4:
			handleDeleteTask()
		case 5:
			return
		default:
			printErrf("Invalid Command\n")
		}
		pressEnterToContinue()
	}
}

func handleCreateTask() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter task title: ")
	scanner.Scan()
	title := scanner.Text()
	fmt.Print("Enter task description: ")
	scanner.Scan()
	description := scanner.Text()
	jsonbody := fmt.Sprintf(`{"title": "%s", "desc": "%s"}`, title, description)

	request, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/tasks",
		strings.NewReader(jsonbody),
	)
	if err != nil {
		printErrf("Failed to create request. err: %s", err.Error())
		return
	}
	request.Header.Set("Authorization", token)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		printErrf("unexpected error while http request.\n%s\n", err.Error())
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusCreated {
		fmt.Println("Task created successfully")
	} else {
		message, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Failed to create task and failed to get error message")
		}
		fmt.Println("Failed to create task. err: ", string(message))
	}
}

func handleShowTask() {
	request, err := http.NewRequest(http.MethodGet, "http://localhost:8080/tasks", nil)
	if err != nil {
		printErrf("Failed to create request. err: %s", err.Error())
		return
	}
	request.Header.Set("Authorization", token)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		printErrf("unexpected error while http request.\n%s\n", err.Error())
		return
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		printErrf("Failed to get tasks and error message")
	}

	if response.StatusCode == http.StatusOK {
		err := json.Unmarshal(data, &tasks)
		if err != nil {
			printErrf("Failed to json tasks\n%s", err.Error())
			os.Exit(1)
		}
		printTasks()
	}
}

func handleUpdateTask() {
	var SNo int
	var id string
	var title string
	var desc string
	var statusInput int
	var status string
	// Scanner := bufio.NewReader(os.Stdin)
	scanner := bufio.NewScanner(os.Stdin)

	handleShowTask()
	fmt.Printf("\nSNo of task to update: ")
	fmt.Scan(&SNo)
	if SNo > len(tasks) {
		printErrf("Invalid SNo\n")
		return
	}
	id = tasks[SNo-1].ID
	fmt.Print("Updated Title (leave empty to not update): ")
	scanner.Scan()
	title = scanner.Text()
	fmt.Print("Updated Description (leave empty to not update): ")
	scanner.Scan()
	desc = scanner.Text()
	fmt.Print("Choose Status.\n1. Pending\n2. Done\n3. Not update\nChoose: ")
	fmt.Scan(&statusInput)

	switch statusInput {
	case 1:
		status = "Pending"
	case 2:
		status = "Done"
	default:
		status = ""
	}

	body := strings.NewReader(
		fmt.Sprintf(`{"title":"%s", "desc": "%s", "status": "%s"}`, title, desc, status),
	)
	request, err := http.NewRequest(http.MethodPut, "http://localhost:8080/tasks/"+id, body)
	if err != nil {
		printErrf("Failed to create request for updating task")
	}
	request.Header.Set("Authorization", token)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		printErrf("unexpected error while http request.\n%s\n", err.Error())
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNoContent {
		fmt.Println("task updated successfully")
		return
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		printErrf("Failed to get tasks and error message")
	}
	printErrf("Failed to update task. err: %s", data)
}

func handleDeleteTask() {
	var SNo int
	var id string

	handleShowTask()
	fmt.Printf("\nSNo of task to delete: ")
	fmt.Scan(&SNo)
	if SNo > len(tasks) {
		printErrf("Invalid SNo\n")
		return
	}
	id = tasks[SNo-1].ID

	request, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/tasks/"+id, nil)
	if err != nil {
		printErrf("Failed to create request for updating task")
	}
	request.Header.Add("Authorization", token)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		printErrf("unexpected error while http request.\n%s\n", err.Error())
		return
	}

	if response.StatusCode == http.StatusNoContent {
		fmt.Println("Task Deleted Successfully")
		return
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		printErrf("Failed to get tasks and error message")
	}
	printErrf("Failed to update task. err: %s", data)
}

func printErrf(s string, a ...any) {
	fmt.Fprintf(os.Stderr, s, a...)
}

func printTasks() {
	if len(tasks) == 0 {
		fmt.Println("no tasks till now")
		return
	}
	fmt.Printf("\n\n")
	titleLeftPadding := (titleWidth - 5) / 2
	title := strings.Repeat(
		" ",
		titleLeftPadding,
	) + "title" + strings.Repeat(
		" ",
		titleWidth-titleLeftPadding-5,
	)
	descLeftPadding := (descWidth - 5) / 2
	desc := strings.Repeat(
		" ",
		descLeftPadding,
	) + "desc" + strings.Repeat(
		" ",
		descWidth-descLeftPadding-4,
	)
	fmt.Printf("SNo|%s|%s|      Status\n", title, desc)
	// fmt.Println("")
	horizontalLine := "---+" + strings.Repeat(
		"-",
		titleWidth,
	) + "+" + strings.Repeat(
		"-",
		descWidth,
	) + "+" + strings.Repeat(
		"-",
		statusWidth,
	)
	fmt.Println(horizontalLine)

	for i, task := range tasks {
		fmt.Printf("%-3d|%-40s|%-60s|%-20s\n", i+1, task.Title, task.Desc, task.Status)
	}
	fmt.Printf("\n\n")
}

# TODO LIST

A task management system built with golang

## Project Setup

- Install dependencies
```bash
go mod tidy
```
- Build api server
```bash
go build -o todo cmd/web/main.go
```
- Run api server
```bash
./todo
```
- Run cli client
```bash
go run tools/client/main.go
```
- Run Tests
```bash
go test -v -coverprofile=cover.out ./internal/...
```
- Get html file for coverage visual
```bash
go tool cover -html=cover.out
```

## Features
- Add task with title, description and status
- Update or delete a task
- List tasks by status
- Save and load task from a local file
- User Based task management

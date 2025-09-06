package models

import (
	"errors"
	"strconv"
)

type TaskRequestDto struct {
	Title  string `json:"title,omitempty"`
	Desc   string `json:"desc,omitempty"`
	Status string `json:"status,omitempty"`
}

func (t TaskRequestDto) ToTask() Task {
	var status int
	switch t.Status {
	case "Pending":
		status = 0
	case "Done":
		status = 1
	default:
		status = 2
	}
	return Task{
		Title:  t.Title,
		Desc:   t.Desc,
		Status: status,
	}
}

type TaskResponseDto struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Desc   string `json:"desc"`
	Status string `json:"status"`
}

func (t TaskResponseDto) ToTask() (Task, error) {
	id, err := strconv.ParseInt(t.ID, 10, 64)
	if err != nil {
		return Task{}, errors.New("invalid Id")
	}

	var status int
	switch t.Status {
	case "Pending":
		status = 0
	case "Done":
		status = 1
	default:
		status = 2
	}

	return Task{
		ID:     id,
		Title:  t.Title,
		Desc:   t.Desc,
		Status: status,
	}, nil
}

package models

type TaskRequestDto struct {
	Title  string `json:"title,omitempty"`
	Desc   string `json:"desc,omitempty"`
	Status string `json:"status,omitempty"`
}

func (trd TaskRequestDto) IsValidStatus() bool {
	switch trd.Status {
	default:
		return false
	case "Pending", "Done", "Waiting":
		return true
	}
}

func (trd TaskRequestDto) ToTask() Task {
	var status int
	switch trd.Status {
	case "Pending":
		status = 0
	case "Done":
		status = 1
	case "Waiting":
		status = 2
	default:
		status = -1
	}
	return Task{
		Title:  trd.Title,
		Desc:   trd.Desc,
		Status: status,
	}
}

type TaskResponseDto struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Desc   string `json:"desc"`
	Status string `json:"status"`
}

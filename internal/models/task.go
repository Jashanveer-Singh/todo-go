package models

import "strconv"

type Task struct {
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	Desc   string `json:"desc"`
	Status int    `json:"status"`
	UserID int64  `json:"user_id"`
}

func (t Task) IsValidTask() bool {
	if t.Title == "" || t.Desc == "" || !t.IsValidStatus() {
		return false
	}

	return true
}

func (t Task) StatusAsText() string {
	switch t.Status {
	case 0:
		return "Pending"
	case 1:
		return "Done"
	default:
		return "Invalid Status"
	}
}

func (t Task) IsValidStatus() bool {
	return t.Status == 0 || t.Status == 1
}

func (t Task) ToDto() TaskResponseDto {
	return TaskResponseDto{
		ID:     strconv.FormatInt(t.ID, 10),
		Title:  t.Title,
		Desc:   t.Desc,
		Status: t.StatusAsText(),
	}
}

package models

import (
	"testing"
)

func TestTask_IsValidTask(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want bool
		task Task
	}{
		{
			name: "task is valid",
			want: true,
			task: Task{
				ID:     123,
				Title:  "title",
				Desc:   "desc",
				Status: 1,
			},
		},
		{
			name: "invalid title",
			want: false,
			task: Task{
				ID:     123,
				Title:  "",
				Desc:   "desc",
				Status: 1,
			},
		},
		{
			name: "invalid desc",
			want: false,
			task: Task{
				ID:     123,
				Title:  "title",
				Desc:   "",
				Status: 1,
			},
		},
		{
			name: "invalid status`",
			want: false,
			task: Task{
				ID:     123,
				Title:  "title",
				Desc:   "desc",
				Status: -91,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.task.IsValidTask()
			if got != tt.want {
				t.Errorf("IsValidTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_StatusAsText(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want string
		task Task
	}{
		{
			name: "status Pending",
			want: "Pending",
			task: Task{
				Status: 0,
			},
		},
		{
			name: "status Done",
			want: "Done",
			task: Task{
				Status: 1,
			},
		},
		{
			name: "status Waiting",
			want: "Waiting",
			task: Task{
				Status: 2,
			},
		},
		{
			name: "status invalid",
			want: "Invalid Status",
			task: Task{
				Status: -1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.task.StatusAsText()
			if tt.want != got {
				t.Errorf("StatusAsText() = %v, want %v.", got, tt.want)
			}
		})
	}
}

func TestTask_IsValidStatus(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		task Task
		want bool
	}{
		{
			name: "status pending",
			task: Task{},
			want: true,
		},
		{
			name: "status is done",
			task: Task{
				Status: 1,
			},
			want: true,
		},
		{
			name: "status is waiting",
			task: Task{
				Status: 2,
			},
			want: true,
		},
		{
			name: "status is invalid",
			task: Task{
				Status: -1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.task.IsValidStatus()
			if tt.want != got {
				t.Errorf("IsValidStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_ToDto(t *testing.T) {
	ta := Task{
		Title:  "task title",
		Desc:   "task desc",
		ID:     12345,
		Status: 1,
	}
	want := TaskResponseDto{
		Title:  "task title",
		Desc:   "task desc",
		ID:     "12345",
		Status: "Done",
	}
	got := ta.ToDto()
	if got != want {
		t.Errorf("ToDto() = %v, want %v", got, want)
	}
}

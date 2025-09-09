package models

import (
	"testing"
)

func TestTaskRequestDto_IsValidStatus(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		taskreq TaskRequestDto
		want    bool
	}{
		{
			name: "status is pending",
			taskreq: TaskRequestDto{
				Status: "Pending",
			},
			want: true,
		},
		{
			name: "status is done",
			taskreq: TaskRequestDto{
				Status: "Done",
			},
			want: true,
		},
		{
			name: "status is Waiting",
			taskreq: TaskRequestDto{
				Status: "Waiting",
			},
			want: true,
		},
		{
			name: "status is invalid",
			taskreq: TaskRequestDto{
				Status: "Pend",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.taskreq.IsValidStatus()
			if got != tt.want {
				t.Errorf("IsValidStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskRequestDto_ToTask(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		taskreq TaskRequestDto
		want    Task
	}{
		{
			name: "task in pending",
			taskreq: TaskRequestDto{
				Title:  "title",
				Desc:   "desc",
				Status: "Pending",
			},
			want: Task{
				Title:  "title",
				Desc:   "desc",
				Status: 0,
			},
		},
		{
			name: "task in done",
			taskreq: TaskRequestDto{
				Title:  "title",
				Desc:   "desc",
				Status: "Done",
			},
			want: Task{
				Title:  "title",
				Desc:   "desc",
				Status: 1,
			},
		},
		{
			name: "task in Waiting",
			taskreq: TaskRequestDto{
				Title:  "title",
				Desc:   "desc",
				Status: "Waiting",
			},
			want: Task{
				Title:  "title",
				Desc:   "desc",
				Status: 2,
			},
		},
		{
			name: "task in invalid status",
			taskreq: TaskRequestDto{
				Title:  "title",
				Desc:   "desc",
				Status: "asdfj",
			},
			want: Task{
				Title:  "title",
				Desc:   "desc",
				Status: -1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.taskreq.ToTask()
			if got != tt.want {
				t.Errorf("ToTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

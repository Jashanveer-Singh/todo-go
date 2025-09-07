package errr

import (
	"net/http"
	"testing"
)

func TestNewUnexpectedError(t *testing.T) {
	message := "unexpected error"
	got := NewUnexpectedError(message)
	want := &AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
	}
	if *got != *want {
		t.Errorf("NewUnexpectedError() = %v, want %v", got, want)
	}
}

func TestNewNotFoundError(t *testing.T) {
	message := "not found"
	got := NewNotFoundError(message)
	want := &AppError{
		Code:    http.StatusNotFound,
		Message: message,
	}
	if *got != *want {
		t.Errorf("NewNotFoundError() = %v, want %v", got, want)
	}
}

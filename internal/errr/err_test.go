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

func TestNewUnauthenticatedError(t *testing.T) {
	message := "unauthenticated error"
	got := NewUnauthenticatedError(message)
	want := &AppError{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
	if *got != *want {
		t.Errorf("NewUnexpectedError() = %v, want %v", got, want)
	}
}

func TestNewBadRequestError(t *testing.T) {
	message := "bad request error"
	got := NewBadRequestError(message)
	want := &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
	if *got != *want {
		t.Errorf("NewUnexpectedError() = %v, want %v", got, want)
	}
}

func TestNewDuplicateError(t *testing.T) {
	message := "duplicate error"
	got := NewDuplicateError(message)
	want := &AppError{
		Code:    http.StatusConflict,
		Message: message,
	}
	if *got != *want {
		t.Errorf("NewUnexpectedError() = %v, want %v", got, want)
	}
}

func TestNewUnauthorizedError(t *testing.T) {
	message := "unauthorized error"
	got := NewUnauthorizedError(message)
	want := &AppError{
		Code:    http.StatusForbidden,
		Message: message,
	}
	if *got != *want {
		t.Errorf("NewUnexpectedError() = %v, want %v", got, want)
	}
}

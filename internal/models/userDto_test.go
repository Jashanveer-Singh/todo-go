package models

import "testing"

func TestUserRequestDto_ToUser(t *testing.T) {
	urd := UserRequestDto{
		Username: "my user",
		Password: "my password",
	}
	want := User{
		ID:       0,
		Username: "my user",
		Password: "my password",
	}
	got := urd.ToUser()
	if got != want {
		t.Errorf("ToUser() = %v, want %v", got, want)
	}
}

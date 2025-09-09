package models

import "testing"

func TestUser_IsValidUser(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		user User
		want bool
	}{
		{
			name: "username empty",
			user: User{
				ID:       0,
				Username: "",
				Password: "qwerqteqetre",
			},
			want: false,
		},
		{
			name: "short password",
			user: User{
				ID:       0,
				Username: "user",
				Password: "qwe",
			},
			want: false,
		},
		{
			name: "valid user",
			user: User{
				ID:       0,
				Username: "user",
				Password: "qwerqteqetre",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.user.IsValidUser()
			if got != tt.want {
				t.Errorf("IsValidUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

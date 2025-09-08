package models

type User struct {
	ID       int64
	Username string
	Password string
}

func (u User) IsValidUser() bool {
	return len(u.Username) > 0 && len(u.Password) > 7
}

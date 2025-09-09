package models

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u User) IsValidUser() bool {
	return len(u.Username) > 0 && len(u.Password) > 7
}

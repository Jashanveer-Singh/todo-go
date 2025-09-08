package models

type UserRequestDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (urd UserRequestDto) ToUser() User {
	return User{
		Username: urd.Username,
		Password: urd.Password,
	}
}

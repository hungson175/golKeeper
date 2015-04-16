package data

import "github.com/hungson175/golKeeper/utils"

type User struct {
	ID       int    `json: "ID"`
	Username string `json: "username"`
	Password string `json: "password"`
}

func (u *User) GetAuthToken() string {
	s := u.Username + u.Password
	return utils.GetMD5Hash(s)
}

func (u *User) equalsExceptID(v *User) bool {
	return (u.Username == v.Username && u.Password == v.Password)
}

type Users []User

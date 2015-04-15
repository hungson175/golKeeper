package data

import "github.com/hungson175/golKeeper/utils"

type User struct {
	id       int
	username string
	password string
}

func (u *User) ID() int {
	return u.id
}

func (u *User) Username() string {
	return u.username
}

func (u *User) Password() string {
	return u.password
}

func (u *User) GetAuthToken() string {
	s := u.username + u.password
	return utils.GetMD5Hash(s)
}

func (u *User) equalsExceptID(v *User) bool {
	return (u.username == v.username && u.password == v.password)
}

type Users []User

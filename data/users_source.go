package data

import (
	"database/sql"
	"errors"
	"fmt"
)

type UsersSource struct {
	db *sql.DB
}

var userDataSource = UsersSource{db: GetDBInstance()}

func GetUserSource() *UsersSource {
	return &userDataSource
}

func CreateUser(username string, password string) (*User, error) {
	inputUser := User{Username: username, Password: password}
	createdUser, err := userDataSource.createUser(&inputUser)
	if err != nil {
		return nil, err
	}
	return createdUser, err
}

type InvalidPasswordError struct{}

func (e InvalidPasswordError) Error() string {
	return "Invalid password"
}

//Login: return user and nil if success, otherwise return not nil error
func Login(username string, password string) (*User, error) {
	ouser, err := userDataSource.getUserByUsername(username)
	if err != nil {
		//TOLEARN: here - in real code, if oyou dont define what error is that => pay later on !
		return nil, err
	}
	if ouser.Password != password {
		return nil, InvalidPasswordError{}
	} else {
		return ouser, nil
	}
}

func ChangePassword(username string, oldPassword string, newPassword string) (*User, error) {
	user, err := Login(username, oldPassword)
	if err != nil {
		return nil, err
	}
	user.Password = newPassword
	user, err = userDataSource.updateUser(user.ID, user)
	if err != nil {
		return nil, err
	}
	return user, err
}

//CRUD

//TODO: test - 2 users with the same user names
//Create: single item in database , return the old user
func (src *UsersSource) createUser(u *User) (*User, error) {
	//TODO: test address of the returned User and the "outside"
	if u == nil {
		panic(errors.New("Nil here user"))
	}
	newUser := *u

	result, err := src.db.Exec("insert into users (username,password) values (?,?)", u.Username, u.Password)
	if err != nil {
		return &newUser, err
	}
	nid, err := result.LastInsertId()
	newUser.ID = int(nid)
	return &newUser, err
}

//Create: batch
func (src *UsersSource) restoreUsersTable(list Users) error {
	db := src.db
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("delete from users")
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, u := range list {
		_, err := tx.Exec(
			"insert into users (id,username,password) values (?,?,?)",
			u.ID, u.Username, u.Password)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

//Read
func (src *UsersSource) getUser(id int) (*User, error) {
	u := User{}
	err := src.db.QueryRow("select * from users where id = ?", id).Scan(&u.ID, &u.Username, &u.Password)
	return &u, err
}

//Read: from username
func (src *UsersSource) getUserByUsername(username string) (*User, error) {
	u := User{}
	err := src.db.QueryRow("select * from users where username=?", username).Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		return nil, err
	}

	return &u, err
}

//Read: all
func (src *UsersSource) getAllUsers() (Users, error) {
	rows, err := src.db.Query("select * from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := Users{}
	for rows.Next() {
		u := User{}
		err := rows.Scan(&u.ID, &u.Username, &u.Password)
		if err != nil {
			return nil, err
		}
		list = append(list, u)
	}
	err = rows.Err()
	return list, err
}

//Update: only change content - not ID !
func (src *UsersSource) updateUser(id int, user *User) (*User, error) {
	result, err := src.db.Exec("update users set username=?, password = ? where id=?", user.Username, user.Password, id)
	if err != nil {
		return nil, err
	}
	rc, _ := result.RowsAffected()
	if rc == 0 {
		return nil, errors.New(fmt.Sprintf("Non-existent id"))
	}
	if rc != 1 {
		return nil, errors.New("Impossible")
	}
	returnUser := *user
	idd, _ := result.LastInsertId()
	returnUser.ID = int(idd)
	return &returnUser, nil
}

//Delete: return
func (src *UsersSource) deleteUser(id int) error {
	_, err := src.db.Exec("delete from users where id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

//Delete: all
func (src *UsersSource) clearData() error {
	_, err := src.db.Exec("delete from users")
	return err
}

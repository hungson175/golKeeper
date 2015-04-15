package data

import (
	"database/sql"
	"errors"
)

type UsersSource struct {
	db *sql.DB
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

	result, err := src.db.Exec("insert into users (username,password) values (?,?)", u.username, u.password)
	if err != nil {
		return &newUser, err
	}
	nid, err := result.LastInsertId()
	newUser.id = int(nid)
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
			u.id, u.username, u.password)
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
	err := src.db.QueryRow("select * from users where id = ?", id).Scan(&u.id, &u.username, &u.password)
	return &u, err
}

//Read: from username
func (src *UsersSource) getUserByUsername(username string) (*User, error) {
	u := User{}
	err := src.db.QueryRow("select * from users where username=?", username).Scan(&u.id, &u.username, &u.password)
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
		err := rows.Scan(&u.id, &u.username, &u.password)
		if err != nil {
			return nil, err
		}
		list = append(list, u)
	}
	err = rows.Err()
	return list, err
}

//Update
func (src *UsersSource) updateUser(id int) (*User, error) {
	return nil, NotImplementedError{"updateUser"}
}

//Delete
func (src *UsersSource) deleteUser(id int) error {
	return NotImplementedError{"deleteUser"}
}

//Delete: all
func (src *UsersSource) clearData() error {
	_, err := src.db.Exec("delete from users")
	return err
}

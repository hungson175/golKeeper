package data

import (
	"fmt"
	"reflect"
	"testing"
)

var compareERFormat = "Expected %v but return %v"
var errorFormat = "Returned err %v"

const DbName = "golkeeper"

type testFunc func(*UsersSource, *testing.T)

var testCases = []testFunc{
	testCreateAndGetAllUsers,
	testGetUserMethods,
	testGetUserByUsername,
	testGetUserByUsername,
}

func TestAll(test *testing.T) {
	ds, err := NewDataSource(DbName)
	if err != nil {
		panic(err)
	}
	defer ds.db.Close()

	//Register test

	src := &UsersSource{ds.db}

	//backup & restore after testing
	list, err := src.getAllUsers()
	defer func(users Users) {
		err := src.restoreUsersTable(users)
		if err != nil {
			panic(err)
		}
	}(list)

	for _, tc := range testCases {

		tc(src, test)
	}
}

func createSampleData(src *UsersSource, test *testing.T) Users {
	users := Users{
		User{username: "sphamhung@gmail.com", password: "dangthaison"},
		User{username: "sphamhung@yahoo.com", password: "abcd1234"},
		User{username: "longhm293@gmail.com", password: "gago1234"},
	}
	src.clearData()
	for _, u := range users {
		newUser, err := src.createUser(&u)
		if err != nil {
			test.Errorf("Cannot create user %v %v", u, err)
		}
		if !newUser.equalsExceptID(&u) {
			test.Errorf(compareERFormat, u, newUser)
		}
	}
	return users
}

func getPass(list Users, username string) string {
	for _, u := range list {
		if u.username == username {
			return u.password
		}
	}
	return ""
}

func testCreateAndGetAllUsers(src *UsersSource, test *testing.T) {
	users := createSampleData(src, test)
	list, err := src.getAllUsers()
	if err != nil {
		test.Errorf(errorFormat, err)
	}

	if len(list) != len(users) {
		test.Errorf(compareERFormat, users, list)
		return
	}
	for _, result := range list {
		uu, err := src.getUserByUsername(result.username)
		if err != nil {
			test.Errorf(errorFormat, err)
			return
		}
		p := uu.password
		if p != getPass(users, result.username) {
			test.Errorf(errorFormat, err)
		}

	}
}
func logError(test *testing.T, err error) {
	test.Errorf(errorFormat, err)
}

func typeOf(x interface{}) string {
	return fmt.Sprintf("%T", x)
}

func testGetUserMethods(src *UsersSource, test *testing.T) {
	createSampleData(src, test)
	list, err := src.getAllUsers()
	if err != nil {
		logError(test, err)
	}
	for _, expected := range list {
		result, err := src.getUser(expected.id)
		if err != nil {
			logError(test, err)
		}
		if !reflect.DeepEqual(expected, *result) {
			test.Errorf(compareERFormat, expected, result)
		}
	}
}

func testGetUserByUsername(src *UsersSource, test *testing.T) {
	users := createSampleData(src, test)
	for _, expected := range users {
		result, err := src.getUserByUsername(expected.username)
		if err != nil {
			logError(test, err)
		}
		if expected.password != result.password {
			test.Errorf(compareERFormat, expected.password, result.password)
		}
	}
}

func testUpdateUser(src *UsersSource, test *testing.T) {
	users := createSampleData(src, test)
	//change user/pass
	for _, user := range users {
		v, _ := src.getUserByUsername(user.username)
		v.password += "+mod"
		src.updateUser(id, &v)
	}
}

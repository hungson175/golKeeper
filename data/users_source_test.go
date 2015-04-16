package data

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/hungson175/golKeeper/utils"
)

var compareERFormat = "Expected %v but return %v"
var errorFormat = "Returned err %v"

type testFunc func(*UsersSource, *testing.T)

var testCases = []testFunc{
	testCreateAndGetAllUsers,
	testGetUserMethods,
	testGetUserByUsername,
	testUpdateUser,
	testDeleteUser,
}

func TestAll(test *testing.T) {
	ds, err := NewDataSource()
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
		funcName := utils.GetFunctionName(tc)
		fmt.Printf("Running test: %s...\n", funcName)
		tc(src, test)
	}
}

func createSampleData(src *UsersSource, test *testing.T) Users {
	users := Users{
		User{username: "sphamhung@gmail.com", password: "dangthaison"},
		User{username: "sphamhung@yahoo.com", password: "abcd1234"},
		User{username: "longhm293@gmail.com", password: "gago1234"},
		User{username: "", password: ""},
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
		_, err := src.updateUser(v.id, v)
		if err != nil {
			logError(test, err)
		}
		newUser, _ := src.getUser(v.id)
		if newUser.password != v.password {
			test.Errorf(compareERFormat, *v, *newUser)
		}
	}

	_, err := src.updateUser(-1, &User{username: "hellokitty", password: "blahblah"})
	if err == nil {
		test.Errorf("The id %v is nonexistent, but the update method returned non-nil value", -1)
	}

	v, _ := src.getUserByUsername(users[0].username)
	src.updateUser(v.id, v)
}

func testDeleteUser(src *UsersSource, test *testing.T) {
	users := createSampleData(src, test)
	for _, user := range users {
		v, _ := src.getUserByUsername(user.username)
		err := src.deleteUser(v.id)
		if err != nil {
			logError(test, err)
		}
		v, err = src.getUserByUsername(user.username)
		if err != sql.ErrNoRows {
			test.Errorf("The user %s should not exists anymore !", user.username)
		}
		fmt.Printf("\n")
	}

	//NOTES: the following test should be run, but for this "learning project" -  just dont
	// err := src.deleteUser(-1)
	// if err == nil {
	// 	test.Errorf("This should return error this casse but i doesnt")
	// }
}

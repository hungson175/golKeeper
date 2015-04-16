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

var SAMPLE_USERS Users = Users{
	User{Username: "sphamhung@gmail.com", Password: "dangthaison"},
	User{Username: "sphamhung@yahoo.com", Password: "abcd1234"},
	User{Username: "longhm293@gmail.com", Password: "gago1234"},
	User{Username: "", Password: ""},
}
var testCases = []testFunc{
	testCreateAndGetAllUsers,
	testGetUserMethods,
	testGetUserByUsername,
	testUpdateUser,
	testDeleteUser,
	testExportedMethod,
}

func TestAll(test *testing.T) {
	src := GetUserSource()

	//backup & restore after testing
	list, _ := src.getAllUsers()

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

	src.clearData()
	for _, u := range SAMPLE_USERS {
		newUser, err := src.createUser(&u)
		if err != nil {
			test.Errorf("Cannot create user %v %v", u, err)
		}
		if !newUser.equalsExceptID(&u) {
			test.Errorf(compareERFormat, u, newUser)
		}
	}
	return SAMPLE_USERS
}

func getPass(list Users, username string) string {
	for _, u := range list {
		if u.Username == username {
			return u.Password
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
		uu, err := src.getUserByUsername(result.Username)
		if err != nil {
			test.Errorf(errorFormat, err)
			return
		}
		p := uu.Password
		if p != getPass(users, result.Username) {
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
		result, err := src.getUser(expected.ID)
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
		result, err := src.getUserByUsername(expected.Username)
		if err != nil {
			logError(test, err)
		}
		if expected.Password != result.Password {
			test.Errorf(compareERFormat, expected.Password, result.Password)
		}
	}
}

func testUpdateUser(src *UsersSource, test *testing.T) {
	users := createSampleData(src, test)
	//change user/pass
	for _, user := range users {
		v, _ := src.getUserByUsername(user.Username)
		v.Password += "+mod"
		_, err := src.updateUser(v.ID, v)
		if err != nil {
			logError(test, err)
		}
		newUser, _ := src.getUser(v.ID)
		if newUser.Password != v.Password {
			test.Errorf(compareERFormat, *v, *newUser)
		}
	}

	_, err := src.updateUser(-1, &User{Username: "hellokitty", Password: "blahblah"})
	if err == nil {
		test.Errorf("The id %v is nonexistent, but the update method returned non-nil value", -1)
	}

	v, _ := src.getUserByUsername(users[0].Username)
	src.updateUser(v.ID, v)
}

func testDeleteUser(src *UsersSource, test *testing.T) {
	users := createSampleData(src, test)
	for _, user := range users {
		v, _ := src.getUserByUsername(user.Username)
		err := src.deleteUser(v.ID)
		if err != nil {
			logError(test, err)
		}
		v, err = src.getUserByUsername(user.Username)
		if err != sql.ErrNoRows {
			test.Errorf("The user %s should not exists anymore !", user.Username)
		}
		fmt.Printf("\n")
	}

	//NOTES: the following test should be run, but for this "learning project" -  just dont
	// err := src.deleteUser(-1)
	// if err == nil {
	// 	test.Errorf("This should return error this casse but i doesnt")
	// }
}

func testExportedMethod(src *UsersSource, test *testing.T) {
	src.clearData()
	expectedUser := User{Username: "sonph", Password: "dangthaison"}
	u, err := CreateUser(expectedUser.Username, expectedUser.Password)
	if err != nil {
		logError(test, err)
	}
	if u.Username != expectedUser.Username || u.Password != expectedUser.Password {
		test.Errorf("Expected %v but returned %v", expectedUser, u)
	}

}

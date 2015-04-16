package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hungson175/golKeeper/data"
	"log"
	"net/http"
)

type UserController struct{}

//CreateAccount: create a new account
//return json: success => newlly created user , failed: error - later on: status code + explicit error (e.g: UsernameExisted)
//TODO: POST & parameters list
func (uc *UserController) CreateAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username, password := vars["username"], vars["password"]
	fmt.Printf("Username = %v , password = %v\n", username, password)
	user, err := data.CreateUser(username, password)
	if err != nil {
		//TODO: how to marshall the AuthenticationToken ?
		json.NewEncoder(w).Encode(err)
		return
	}
	if err := json.NewEncoder(w).Encode(*user); err != nil {
		log.Fatal(err)
		panic(err)
	}
}

//Login
func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username, password := vars["username"], vars["password"]
	user, err := data.Login(username, password)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	if err := json.NewEncoder(w).Encode(*user); err != nil {
		log.Fatal(err)
		panic(err)
	}
}

//Change password
func (uc *UserController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username, oldPass, newPass := vars["username"], vars["oldpassword"], vars["newpassword"]
	user, err := data.ChangePassword(username, oldPass, newPass)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	if err := json.NewEncoder(w).Encode(*user); err != nil {
		log.Fatal(err)
		panic(err)
	}

}

package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/hungson175/golKeeper/data"
)

type Page struct {
	Title    string
	Username string
}

func Signup(w http.ResponseWriter, r *http.Request) {
	type suData struct {
		Username string
	}
	session, err := store.Get(r, SESSION_NAME)
	if err != nil {
		log.Fatal(err)
		return
	}
	username := session.Values[SSK_USER_NAME].(string)
	data := &suData{Username: username}
	renderTemplate(w, "signup.html", data)
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	cfPassword := r.FormValue("confirmed_password") //TODO: this should be done on Client, not server
	if password != cfPassword {
		http.Redirect(w, r, "/Signup", http.StatusFound)
		return
	}
	_, err := data.CreateUser(username, password)
	if err != nil {
		log.Fatal(err)
		http.Redirect(w, r, "/Signup", http.StatusInternalServerError)
		return
	}
	// fmt.Println("Done creating user")
	session, err := store.Get(r, SESSION_NAME)
	session.Values[SSK_USER_NAME] = username
	err = session.Save(r, w)
	if err == nil {
		http.Redirect(w, r, "/PersonalPage", http.StatusFound)
		// fmt.Println("Done redirecct to personal page")
	} else {
		http.Redirect(w, r, "/Signup", http.StatusInternalServerError)
		// fmt.Println("Done redirecct to sign up page")
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "Login"}
	session, err := store.Get(r, SESSION_NAME)
	if err != nil {
		log.Fatal(err)
	} else {
		username := session.Values[SSK_USER_NAME]
		if username != nil {
			p.Username = username.(string)
		}
	}

	renderTemplate(w, "login.html", p)
}

const SESSION_NAME = "user-session"
const USER_AUTH_TOKEN = "user_token"
const SSK_USER_ID = "userID"
const SECRECT_KEY = "asda-341h-sd45763"
const SSK_USER_NAME = "ssk_username"
const SSK_PASSWORD = "ssk_password"

var store = sessions.NewCookieStore([]byte(SECRECT_KEY)) //TODO: should be enviroment variable ? No, the deployment will be complicated

func AccountAction(w http.ResponseWriter, r *http.Request) {
	if isLogin := r.FormValue("login_button"); isLogin != "" {
		VerifyAccount(w, r)
	} else if isSignup := r.FormValue("signup_button"); isSignup != "" {
		Signup(w, r)
	} else {
		panic(errors.New("Not a supported form"))
	}
}

func VerifyAccount(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	user, err := data.Login(username, password)
	session, _ := store.Get(r, SESSION_NAME)
	if err == nil {
		session.Values[SSK_USER_ID] = user.ID
		session.Values[SSK_USER_NAME] = user.Username
		err = session.Save(r, w)
		http.Redirect(w, r, "/PersonalPage/", http.StatusFound)
		return
	} else {
		session.Values[SSK_USER_NAME] = username
		session.Save(r, w)
		http.Redirect(w, r, "/Login/", http.StatusFound)
		return
	}
}

package controllers

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/gorilla/sessions"
	"github.com/hungson175/golKeeper/data"
)

type Page struct {
	Title    string
	Username string
}

var (
	doOnce              sync.Once
	userActionTemplates *template.Template
	BASE_DIR            = os.Getenv("GOL_KEEPER_PATH")
)

func renderTemplate(w http.ResponseWriter, controllerName string, p interface{}) {
	err := userActionTemplates.ExecuteTemplate(w, controllerName+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func init() {
	doOnce.Do(

		func() {
			dir := path.Join(BASE_DIR, "templates")
			userActionTemplates = template.Must(template.ParseFiles(
				dir+"/login.html",
				dir+"/signup.html"))
		})
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
	renderTemplate(w, "signup", data)
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	cfPassword := r.FormValue("confirmed_password") //TODO: this should be done on Client, not server
	fmt.Printf("username = %s , password = %s\n", username, password)
	if password != cfPassword {
		fmt.Println("Redirect to signup")
		http.Redirect(w, r, "/Signup", http.StatusFound)
		return
	}
	fmt.Println("Creating user")
	_, err := data.CreateUser(username, password)
	if err != nil {
		log.Fatal(err)
		http.Redirect(w, r, "/Signup", http.StatusInternalServerError)
		return
	}
	fmt.Println("Done creating user")
	session, err := store.Get(r, SESSION_NAME)
	session.Values[SSK_USER_NAME] = username
	err = session.Save(r, w)
	if err == nil {
		http.Redirect(w, r, "/PersonalPage", http.StatusFound)
		fmt.Println("Done redirecct to personal page")
	} else {
		http.Redirect(w, r, "/Signup", http.StatusInternalServerError)
		fmt.Println("Done redirecct to sign up page")
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

	renderTemplate(w, "login", p)
}

const SESSION_NAME = "user-session"
const USER_AUTH_TOKEN = "user_token"
const USER_ID_SESSION_KEY = "userID"
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
		session.Values[USER_ID_SESSION_KEY] = user.ID
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

func PersonalPage(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, SESSION_NAME)
	if err != nil {
		http.Redirect(w, r, "/Login", http.StatusBadRequest)
		return
	}
	username := session.Values[SSK_USER_NAME]
	fmt.Fprintf(w, "Hello %s", username)
}

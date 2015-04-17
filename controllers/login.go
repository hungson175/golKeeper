package controllers

import (
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

func renderTemplate(w http.ResponseWriter, controllerName string, p *Page) {
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
				dir + "/login.html"))
		})
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

func VerifyAccount(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside VerifyAccount()")
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Printf("Username = %s and Password = %s", username, password)
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

func Debug(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside Debug()")
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Printf("Username = %s and Password = %s", username, password)
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

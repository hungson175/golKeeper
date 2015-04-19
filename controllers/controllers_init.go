package controllers

import (
	"html/template"
	"net/http"
	"os"
	"path"
	"sync"
)

var (
	doOnce    sync.Once
	templates *template.Template
	BASE_DIR  = os.Getenv("GOL_KEEPER_PATH")
)

const (
	FV_DESC         = "desc"
	FV_DATE         = "start_date"
	FV_WEEKS        = "weeks"
	FV_REF_USERNAME = "ref_name"
	FV_LOSE_COMMIT  = "lose_commit"
)

func init() {
	doOnce.Do(
		func() {
			dir := path.Join(BASE_DIR, "templates")
			templates = template.Must(template.ParseFiles(
				dir+"/login.html",
				dir+"/signup.html",
				dir+"/edit_goal.html",
				dir+"/personal.html"))
		})
}

func renderTemplate(w http.ResponseWriter, templateName string, p interface{}) {
	err := templates.ExecuteTemplate(w, templateName, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/Login", http.StatusUnauthorized)
}

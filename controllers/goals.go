package controllers

import (
	"net/http"
	"strconv"

	"github.com/hungson175/golKeeper/data"
	"github.com/hungson175/golKeeper/utils"
)

func PersonalPage(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, SESSION_NAME)
	username := session.Values[SSK_USER_NAME]
	type tmpData struct {
		Username string
	}
	data := tmpData{Username: username.(string)}
	renderTemplate(w, "personal.html", data)
}

func EditGoal(w http.ResponseWriter, r *http.Request) {
	type tmpData struct {
		Username string
	}
	session, err := store.Get(r, SESSION_NAME)
	if err != nil {
		redirectToLogin(w, r)
		return
	}
	username := session.Values[SSK_USER_NAME].(string)
	data := tmpData{Username: username}
	renderTemplate(w, "edit_goal.html", data)
}

func SaveGoal(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, SESSION_NAME)
	if err != nil {
		redirectToLogin(w, r)
		return
	}
	goal := &data.Goal{}
	goal.Desc = r.FormValue(FV_DESC)
	goal.StartDate = utils.String2Date(r.FormValue(FV_DATE))
	goal.Weeks, _ = strconv.Atoi(r.FormValue(FV_WEEKS))
	goal.UserID = session.Values[SSK_USER_ID].(int)
	goal.RefName = r.FormValue(FV_REF_USERNAME)
	//RefID
	refUser := data.GetUserByUsername(goal.RefName)
	if refUser == nil {
		http.Redirect(w, r, "/EditGoal", http.StatusBadRequest)
		return
	}
	goal.RefID = refUser.ID
	goal.LoseCommit = r.FormValue(FV_LOSE_COMMIT)
	_, err = data.CreateGoal(goal)
	if err != nil {
		http.Redirect(w, r, "/EditGoal", http.StatusInternalServerError)
		return
	} else {
		http.Redirect(w, r, "/PersonalPage", http.StatusFound)
	}
}

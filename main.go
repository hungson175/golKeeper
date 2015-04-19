package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hungson175/golKeeper/controllers"
	"github.com/hungson175/golKeeper/controllers/apis"
)

func main() {
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}
type Routes []Route

var userAPIRoutes = Routes{
	//TODO: how to post this info ?
	Route{"Create Account", "GET", "/create/{username}/{password}", apis.CreateAccount},
	Route{"Login", "GET", "/login/{username}/{password}", apis.Login},
	Route{"Change password", "GET", "/changepassword/{username}/{oldpassword}/{newpassword}", apis.ChangePassword},
}

var webRoutes = Routes{
	Route{"Login page", "GET", "/Login", controllers.Login},
	Route{"Verify account page", "POST", "/VerifyAccount", controllers.VerifyAccount},
	Route{"Signup", "GET", "/Signup", controllers.Signup},
	Route{"Select Account Action", "POST", "/AccountAction", controllers.AccountAction},
	Route{"Create account", "POST", "/CreateAccount", controllers.CreateAccount},
	Route{"Personal page", "GET", "/PersonalPage", controllers.PersonalPage},
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	// apiRouter := router.PathPrefix("/apis").Subrouter()
	// registerUserAPIRoutes(apiRouter.PathPrefix("/user").Subrouter(), userAPIRoutes)
	registerUserAPIRoutes(router, webRoutes)
	return router
}

func registerUserAPIRoutes(router *mux.Router, routes Routes) {
	//register routes for user api
	for _, route := range routes {
		handler := Logger(route.HandlerFunc, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

}

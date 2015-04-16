package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hungson175/golKeeper/controllers"
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

var userController controllers.UserController = controllers.UserController{}
var userAPIRoutes = Routes{
	//TODO: how to post this info ?
	Route{"Create Account", "GET", "/create/{username}/{password}", userController.CreateAccount},
	Route{"Login", "GET", "/login/{username}/{password}", userController.Login},
	Route{"Change password", "GET", "/changepassword/{username}/{oldpassword}/{newpassword}", userController.ChangePassword},
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	apiRouter := router.PathPrefix("/apis").Subrouter()
	registerUserAPIRoutes(apiRouter, userAPIRoutes, "user")
	return router
}

func registerUserAPIRoutes(router *mux.Router, routes Routes, controllerName string) {
	subRouter := router.PathPrefix("/" + controllerName).Subrouter()
	//register routes for user api
	for _, route := range routes {
		handler := Logger(route.HandlerFunc, route.Name)
		subRouter.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

}

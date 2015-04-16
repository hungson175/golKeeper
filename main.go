package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
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
userAPIRoutes := Routes {
  //TODO: how to post this info ?
  Route{"Create Account","GET","/create/{username}/{pasword}",userController.CreateAccount},
  Route{"Login","GET","/login/{username}/{pasword}",userController.Login},
  Route{"Change password","GET","/changepassword/{username}/{oldpasword}/{newpassword}",userController.ChangePassword},
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	apiRouter := router.PathPrefix("/apis").Subrouter()
}

func registerUserAPIRoutes() {
  userAPIRouter := apiRouter.PathPrefix("/user")
  //register routes for user api
  for _, route := range userAPIRoutes {
    handler := Logger(route.HandlerFunc, route.Name)
    userAPIRouter.
      Methods(route.Method).
      Path(route.Pattern).
      Name(route.Name).
      Handler(handler)
  }

}

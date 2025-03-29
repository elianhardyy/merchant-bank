package routes

import (
	"go-json/internal/controllers"
	"go-json/internal/middlewares"
	"go-json/internal/security"
	"net/http"
)

func UserRoutes(api controllers.UserController, token security.TokenService) {
	auth := R.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", api.Register).Methods("POST")
	auth.HandleFunc("/login", api.Login).Methods("POST")
	auth.Handle("/logout", middlewares.ProtectedHandler(http.HandlerFunc(api.Logout), token, []string{"customer", "merchant"})).Methods("POST")
	// auth.Handle("/refresh", middlewares.ProtectedHandler(http.HandlerFunc(api.Refresh), token, []string{"customer", "merchant"})).Methods("POST")
	user := R.PathPrefix("/user").Subrouter()
	user.Handle("/users", middlewares.ProtectedHandler(http.HandlerFunc(api.UserList), token, []string{"merchant"})).Methods("GET")
}

package routes

import (
	"go-json/internal/controllers"
	"go-json/internal/middlewares"
	"go-json/internal/security"
	"net/http"
)

func TransactionRoutes(api controllers.TransactionController, token security.TokenService) {
	transaction := R.PathPrefix("/trx").Subrouter()
	transaction.Handle("/create", middlewares.ProtectedHandler(http.HandlerFunc(api.Payment), token, []string{"customer"})).Methods("POST")
	transaction.Handle("/history/{id}", middlewares.ProtectedHandler(http.HandlerFunc(api.TransactionHistory), token, []string{"customer", "merchant"})).Methods("GET")
}

package router

import (
	"github.com/gorilla/mux"

	"main/backend/handlers"
	"main/backend/utils"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	// Middleware
	r.Use(utils.LoggingMiddleware)
	r.Use(utils.AuthMiddleware)

	r.HandleFunc("/goodstocks", handlers.GetGoodStocksHandler).Methods("GET")
	r.HandleFunc("/ticker/{id:[a-zA-Z]+}", handlers.GetLiveStockData).Methods("GET")
	// r.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	// r.HandleFunc("/users/{id}", handlers.GetUser).Methods("GET")
	// r.HandleFunc("/products", handlers.GetProducts).Methods("GET")
	// r.HandleFunc("/products/{id}", handlers.GetProduct).Methods("GET")

	return r
}

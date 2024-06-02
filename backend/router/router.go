package router

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"main/backend/handlers"
	"main/backend/utils"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	// Middleware
	r.Use(utils.LoggingMiddleware)
	r.Use(utils.AuthMiddleware)

	r.HandleFunc("/", handlers.IndexHandler).Methods("GET")
	r.HandleFunc("/goodstocks", handlers.GetGoodStocksHandler).Methods("GET")
	r.HandleFunc("/ticker/{id:[a-zA-Z]+}", handlers.GetLiveStockData).Methods("GET")
	r.Handle("/metrics", promhttp.Handler())
	return r
}

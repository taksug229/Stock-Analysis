package router

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"main/backend/handlers"
	"main/backend/utils"
)

func NewRouter() {
	r := mux.NewRouter()

	// Middleware
	r.Use(utils.LoggingMiddleware)
	r.Use(utils.AuthMiddleware)

	r.HandleFunc("/", handlers.IndexHandler).Methods("GET")
	r.HandleFunc("/goodstocks", handlers.GetGoodStocksHandler).Methods("GET")
	r.HandleFunc("/ticker/{id:[a-zA-Z]+}", handlers.GetLiveStockData).Methods("GET")
	// r.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	// r.HandleFunc("/users/{id}", handlers.GetUser).Methods("GET")
	// r.HandleFunc("/products", handlers.GetProducts).Methods("GET")
	// r.HandleFunc("/products/{id}", handlers.GetProduct).Methods("GET")
	r.Handle("/metrics", promhttp.Handler())
	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	log.Println("Starting server on :8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}

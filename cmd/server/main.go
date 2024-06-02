// cmd/server/main.go
package main

import (
	"log"
	"main/backend/router"
	"net/http"
	"time"
)

func main() {
	log.Println("Starting application server...")
	r := router.NewRouter()
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

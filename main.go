package main

import (
	"log"
	"net/http"
	"time"

	// "main/api"
	"main/router"
)

func main() {
	// uniqueTickers := api.SaveFinancialData()
	// uniqueTickers["VOO"] = struct{}{}
	// api.SaveQuoteFromYahoo(uniqueTickers)
	// gcp.UploadToGCSToBigQuery()
	// gcp.CreateMLTable()

	// router.HandleFunc("/resources/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "/resources/%s!", mux.Vars(r)["id"])
	// })
	// router.HandleFunc("/resources/{id:[0-9]+}/values", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "/resources/values: %s", mux.Vars(r)["id"])
	// })
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

	log.Println("Completed!")
}

package main

import (
	"main/backend/api"
	"main/backend/gcp"
	// "time"
	// "main/backend/api"
	// "main/backend/gcp"
	// "main/backend/router"
)

func main() {
	uniqueTickers := api.SaveFinancialData()
	uniqueTickers["VOO"] = struct{}{}
	// api.SaveQuoteFromYahoo(uniqueTickers)
	gcp.UploadToGCSToBigQuery()
	// gcp.CreateMLTable()
	// r := router.NewRouter()
	// r.Handle("/metrics", promhttp.Handler())
	// srv := &http.Server{
	// 	Handler:      r,
	// 	Addr:         ":8080",
	// 	WriteTimeout: 15 * time.Second,
	// 	ReadTimeout:  15 * time.Second,
	// 	IdleTimeout:  60 * time.Second,
	// }

	// log.Println("Starting server on :8080")
	// if err := srv.ListenAndServe(); err != nil {
	// 	log.Fatalf("Could not start server: %s\n", err)
	// }
}

// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"time"

// 	"github.com/gorilla/mux"
// )

// func main() {
// 	router := mux.NewRouter()
// 	router.HandleFunc("/resources", func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprintf(w, "/resources")
// 	})
// 	router.HandleFunc("/resources/{id:[a-zA-Z]+}", func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprintf(w, "/resources/: %s!", mux.Vars(r)["id"])
// 	})
// 	router.HandleFunc("/resources/{id:[0-9]+}/values", func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprintf(w, "/resources/{id:[0-9]+}/values: %s", mux.Vars(r)["id"])
// 	})
// 	srv := &http.Server{
// 		Handler:      router,
// 		Addr:         ":8080",
// 		WriteTimeout: 15 * time.Second,
// 		ReadTimeout:  15 * time.Second,
// 	}

// 	log.Fatal(srv.ListenAndServe())
// }

// package main

// import (
// 	"fmt"

// 	"main/backend/handlers"
// 	"net/http"
// )

// func main() {
// 	http.HandleFunc("/", handlers.GetLiveStockData)
// 	fmt.Println("Starting server at http://localhost:8080")
// 	// Start the server on port 8080
// 	err := http.ListenAndServe(":8080", nil)
// 	if err != nil {
// 		fmt.Println("Error starting server:", err)
// 	}
// }

package main

import (
	"main/backend/router"
	// "time"
)

func main() {
	// uniqueTickers := api.SaveFinancialData()
	// api.SaveQuoteFromYahoo(uniqueTickers)
	// gcp.UploadToGCSToBigQuery()
	// gcp.CreateMLTable()
	// gcp.CreateTrainTestTable()
	// gcp.CreateModel()
	router.NewRouter()
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

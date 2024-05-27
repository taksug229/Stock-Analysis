// package main

// import (
// 	"fmt"
// 	"net/http"

// 	// "time"

// 	// "main/backend/api"
// 	// "main/backend/gcp"
// 	// "main/backend/router"

// 	"github.com/gorilla/mux"
// )

// func main() {
// 	// uniqueTickers := api.SaveFinancialData()
// 	// uniqueTickers["VOO"] = struct{}{}
// 	// api.SaveQuoteFromYahoo(uniqueTickers)
// 	// gcp.UploadToGCSToBigQuery()
// 	// gcp.CreateMLTable()
// 	// r := router.NewRouter()
// 	// r.Handle("/metrics", promhttp.Handler())
// 	// srv := &http.Server{
// 	// 	Handler:      r,
// 	// 	Addr:         ":8080",
// 	// 	WriteTimeout: 15 * time.Second,
// 	// 	ReadTimeout:  15 * time.Second,
// 	// 	IdleTimeout:  60 * time.Second,
// 	// }

// 	// 	log.Println("Starting server on :8080")
// 	// 	if err := srv.ListenAndServe(); err != nil {
// 	// 		log.Fatalf("Could not start server: %s\n", err)
// 	// 	}
// }

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

package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// Struct to hold template data
type StockData struct {
	Ticker              string
	CurrentStockPrice   float64
	PredictedStockPrice float64
	MarketValue         float64
	IntrinsicValue      float64
	Recommendation      string
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Create an instance of StockData with the data you want to display
	data := StockData{
		Ticker:              "GOOGL",
		CurrentStockPrice:   1500.50,
		PredictedStockPrice: 1600.00,
		MarketValue:         1_000_000_000.00,
		IntrinsicValue:      1_200_000_000.00,
		Recommendation:      "Buy",
	}

	// Parse the template file
	tmpl, err := template.ParseFiles("frontend/templates/template.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	// Execute the template with the data
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Starting server at http://localhost:8080")
	// Start the server on port 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
